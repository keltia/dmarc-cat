package main

import (
	"bytes"
	"net"
	"text/template"
	"time"

	"github.com/intel/tfortools"
	"github.com/pkg/errors"
)

const (
	reportTmpl = `{{.MyName}} {{.MyVersion}} by {{.Author}}

Reporting by: {{.Org}} — {{.Email}}
From {{.DateBegin}} to {{.DateEnd}}

Domain: {{.Domain}}
Policy: p={{.Disposition}}; dkim={{.DKIM}}; spf={{.SPF}}

Reports({{.Count}}):
`

	rowTmpl = `{{ table (sort . "Count" "dsc")}}`
)

type headVars struct {
	MyName      string
	MyVersion   string
	Author      string
	Org         string
	Email       string
	DateBegin   string
	DateEnd     string
	Domain      string
	Disposition string
	DKIM        string
	SPF         string
	Pct         int
	Count       int
}

type row struct {
	IP    net.IP
	Count int
	From  string
	RFrom string
	RDKIM string
	RSPF  string
}

// GatherRows extracts all IP and return the rows
func GatherRows(r Feedback) []row {
	var rows []row

	for _, report := range r.Records {
		current := row{
			IP:    report.Row.SourceIP,
			Count: report.Row.Count,
			From:  report.Identifiers.HeaderFrom,
		}
		if report.AuthResults.DKIM.Domain == "" {
			current.RFrom = report.AuthResults.SPF.Domain
		} else {
			current.RFrom = report.AuthResults.DKIM.Domain
		}
		current.RSPF = report.AuthResults.SPF.Result
		current.RDKIM = report.AuthResults.DKIM.Result

		rows = append(rows, current)
	}
	return rows
}

// Analyze extract and display what we want
func Analyze(r Feedback) (string, error) {
	var buf bytes.Buffer

	tmplvars := &headVars{
		MyName:      MyName,
		MyVersion:   MyVersion,
		Author:      Author,
		Org:         r.Metadata.OrgName,
		Email:       r.Metadata.Email,
		DateBegin:   time.Unix(r.Metadata.Date.Begin, 0).String(),
		DateEnd:     time.Unix(r.Metadata.Date.End, 0).String(),
		Domain:      r.Policy.Domain,
		Disposition: r.Policy.P,
		DKIM:        r.Policy.ADKIM,
		SPF:         r.Policy.ASPF,
		Pct:         r.Policy.Pct,
		Count:       len(r.Records),
	}

	rows := GatherRows(r)

	// Header
	t := template.Must(template.New("r").Parse(string(reportTmpl)))
	err := t.ExecuteTemplate(&buf, "r", tmplvars)
	if err != nil {
		return "", errors.Wrapf(err, "error in template 'r'")
	}

	err = tfortools.OutputToTemplate(&buf, "reports", rowTmpl, rows, nil)
	if err != nil {
		return "", errors.Wrapf(err, "error in template 'reports'")
	}

	return buf.String(), err
}
