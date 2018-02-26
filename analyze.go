package main

import (
	"bytes"
	"github.com/intel/tfortools"
	"log"
	"net"
	"text/template"
	"time"
)

const (
	reportTmpl = `{{.MyName}} {{.MyVersion}} by {{.Author}}

Reporting by: {{.Org}} — {{.Email}}
From {{.DateBegin}} to {{.DateEnd}}

Domain: {{.Domain}}
Policy: p={{.Disposition}}; dkim={{.DKIM}}; spf={{.SPF}}

Reports({{.Count}}):
`

	rowTmpl = `{{ table .}}`
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
	t := template.Must(template.New("r").Parse(string(reportTmpl)))
	err := t.ExecuteTemplate(&buf, "r", tmplvars)
	if err != nil {
		log.Printf("error in template 'r': %v", err)
	}

	err = tfortools.OutputToTemplate(&buf, "reports", rowTmpl, rows, nil)
	if err != nil {
		log.Printf("error in template 'reports': %v", err)
	}

	return buf.String(), err
}
