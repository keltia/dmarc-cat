package main

import (
	"bytes"
	"fmt"
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

// My template vars
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

// Single row
type row struct {
	IP    string
	Count int
	From  string
	RFrom string
	RDKIM string
	RSPF  string
}

func ResolveIP(ctx *Context, ip string) string {
	ips, err := ctx.r.LookupAddr(ip)
	if err != nil {
		return ip
	}
	// XXX FIXME?
	return ips[0]
}

// GatherRows extracts all IP and return the rows
func GatherRows(ctx *Context, r Feedback) []row {
	var rows []row

	for _, report := range r.Records {

		ip0 := ResolveIP(ctx, report.Row.SourceIP.String())
		current := row{
			IP:    ip0,
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
func Analyze(ctx *Context, r Feedback) (string, error) {
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

	rows := GatherRows(ctx, r)
	if len(rows) == 0 {
		return "", fmt.Errorf("empty report")
	}

	// Header
	t := template.Must(template.New("r").Parse(string(reportTmpl)))
	err := t.ExecuteTemplate(&buf, "r", tmplvars)
	if err != nil {
		return "", errors.Wrapf(err, "error in template 'r'")
	}

	// Rows
	err = tfortools.OutputToTemplate(&buf, "reports", rowTmpl, rows, nil)
	if err != nil {
		return "", errors.Wrapf(err, "error in template 'reports'")
	}

	return buf.String(), nil
}
