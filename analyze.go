package main

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/intel/tfortools"
	"github.com/ivpusic/grpool"
	"github.com/pkg/errors"
)

const (
	reportTmpl = `{{.MyName}} {{.MyVersion}}/j{{.Jobs}} by {{.Author}}

Reporting by: {{.Org}} — {{.Email}}
From {{.DateBegin}} to {{.DateEnd}}

Domain: {{.Domain}}
Policy: p={{.Disposition}}; dkim={{.DKIM}}; spf={{.SPF}}

Reports({{.Count}}):
`

	rowTmpl = `{{ table (sort . %s)}}`
)

// My template vars
type headVars struct {
	MyName      string
	MyVersion   string
	Jobs        string
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

// Entry representes a single entry
type Entry struct {
	IP    string
	Count int
	From  string
	RFrom string
	RDKIM string
	RSPF  string
}

// ResolveIP returns the first name associated with the IP
func ResolveIP(ctx *Context, ip string) string {
	ips, err := ctx.r.LookupAddr(ip)
	if err != nil {
		return ip
	}
	// XXX FIXME?
	return ips[0]
}

type IP struct {
	IP   string
	Name string
}

func ParallelSolve(ctx *Context, iplist []IP) []IP {
	resolved := make([]IP, len(iplist))
	pool := grpool.NewPool(ctx.jobs, len(iplist))
	verbose("ParallelSolve with %d workers", ctx.jobs)
	defer pool.Release()

	pool.WaitCount(len(iplist))
	for i, e := range iplist {
		current := e.IP
		ind := i
		pool.JobQueue <- func() {
			defer pool.JobDone()

			resolved[ind].Name = ResolveIP(ctx, current)
		}
	}
	pool.WaitAll()

	return resolved
}

// GatherRows extracts all IP and return the rows
func GatherRows(ctx *Context, r Feedback) []Entry {
	var (
		rows    []Entry
		iplist  []IP
		newlist []IP
	)

	ipslen := len(r.Records)

	if !fNoResolv {
		verbose("Resolving all %d IPs", ipslen)
		iplist = make([]IP, ipslen)
		// Get all IPs
		for i, report := range r.Records {
			iplist[i] = IP{IP: report.Row.SourceIP.String()}
		}

		// Now we have a nice array
		newlist = ParallelSolve(ctx, iplist)
		verbose("Resolved %d IPs", ipslen)
	} else {
		newlist = make([]IP, ipslen)
		// Get all IPs
		for i, report := range r.Records {
			newlist[i] = IP{Name: report.Row.SourceIP.String()}
		}
	}

	for i, report := range r.Records {
		ip0 := newlist[i].Name
		current := Entry{
			IP:    ip0,
			Count: report.Row.Count,
			From:  report.Identifiers.HeaderFrom,
			RSPF:  report.AuthResults.SPF.Result,
			RDKIM: report.AuthResults.DKIM.Result,
		}
		if report.AuthResults.DKIM.Domain == "" {
			current.RFrom = report.AuthResults.SPF.Domain
		} else {
			current.RFrom = report.AuthResults.DKIM.Domain
		}
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
		Jobs:        fmt.Sprintf("%d", fJobs),
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

	// Generate our template
	sortTmpl := fmt.Sprintf(rowTmpl, fSort)
	err = tfortools.OutputToTemplate(&buf, "reports", sortTmpl, rows, nil)
	if err != nil {
		return "", errors.Wrapf(err, "error in template 'reports'")
	}

	return buf.String(), nil
}
