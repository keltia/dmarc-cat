package main

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/intel/tfortools"
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

type IP struct {
	IP   string
	Name string
}

// ParallelSolve is doing the IP to name resolution with a worker set.
// XXX use Mutex
func ParallelSolve(ctx *Context, iplist []IP) []IP {
	verbose("ParallelSolve with %d workers", ctx.jobs)

	var lock sync.Mutex

	wg := &sync.WaitGroup{}
	queue := make(chan IP, ctx.jobs)

	resolved := iplist

	ind := 0

	for i := 0; i < ctx.jobs; i++ {
		wg.Add(1)

		debug("setting up w%d", i)
		go func(n int, wg *sync.WaitGroup) {
			defer wg.Done()

			var name string

			for e := range queue {
				ips, err := ctx.r.LookupAddr(e.IP)
				debug("ips=%#v", ips)
				if err != nil {
					name = e.IP
				} else {
					name = ips[0]
				}

				lock.Lock()
				resolved[ind].Name = name
				ind++
				lock.Unlock()
				debug("w%d - ip=%s - name=%s", n, e.IP, name)
			}
		}(i, wg)
	}

	for _, q := range iplist {
		queue <- q
	}

	close(queue)
	wg.Wait()

	debug("resolved=%#v", resolved)
	return resolved
}

// GatherRows extracts all IP and return the rows
func GatherRows(ctx *Context, r Report) []Entry {
	var (
		rows    []Entry
		iplist  []IP
		newlist []IP
	)

	ipslen := len(r.Records)

	verbose("Resolving all %d IPs", ipslen)
	iplist = make([]IP, ipslen)
	// Get all IPs
	for i, report := range r.Records {
		iplist[i] = IP{IP: report.Row.SourceIP.String()}
	}

	// Now we have a nice array
	newlist = ParallelSolve(ctx, iplist)
	verbose("Resolved %d IPs", ipslen)

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

// TODO: fix Analyze to work with multiple Reports.

// Analyze extract and display what we want
// XXX only analyse the first Report in a feedback.
func Analyze(ctx *Context, r Feedback) (string, error) {
	var buf bytes.Buffer

	single := r[0]
	tmplvars := &headVars{
		MyName:      MyName,
		MyVersion:   MyVersion,
		Jobs:        fmt.Sprintf("%d", fJobs),
		Author:      Author,
		Org:         single.Metadata.OrgName,
		Email:       single.Metadata.Email,
		DateBegin:   time.Unix(single.Metadata.Date.Begin, 0).String(),
		DateEnd:     time.Unix(single.Metadata.Date.End, 0).String(),
		Domain:      single.Policy.Domain,
		Disposition: single.Policy.P,
		DKIM:        single.Policy.ADKIM,
		SPF:         single.Policy.ASPF,
		Pct:         single.Policy.Pct,
		Count:       len(single.Records),
	}

	rows := GatherRows(ctx, single)
	if len(rows) == 0 {
		return "", fmt.Errorf("empty report")
	}

	// Header
	t := template.Must(template.New("r").Parse(reportTmpl))
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
