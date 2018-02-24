package main

import (
	"net"
)

type DateRange struct {
	Begin int64 `xml:"begin"`
	End   int64 `xml:"end"`
}

type ReportMetadata struct {
	OrgName          string    `xml:"org_name"`
	Email            string    `xml:"email"`
	ExtraContactInfo string    `xml:"extra_contact_info"`
	ReportID         string    `xml:"report_id"`
	Date             DateRange `xml:"date_range"`
}

type PolicyPublished struct {
	Domain string `xml:"domain"`
	ADKIM  string `xml:"adkim"`
	ASPF   string `xml:"aspf"`
	P      string `xml:"p"`
	SP     string `xml:"sp"`
	Pct    int    `xml:"pct"`
}

type PolicyEvaluated struct {
	Disposition string `xml:"disposition"`
	DKIM        string `xml:"dkim"`
	SPF         string `xml:"spf"`
}

type Row struct {
	SourceIP net.IP          `xml:"source_ip"`
	Count    int             `xml:"count"`
	Policy   PolicyEvaluated `xml:"policy_evaluated"`
}

type Identifiers struct {
	HeaderFrom string `xml:"header_from"`
}

type Result struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"`
}

type AuthResults struct {
	DKIM Result `xml:"dkim"`
	SPF  Result `xml:"spf"`
}

type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResults `xml:"auth_results"`
}

type Feedback struct {
	Metadata ReportMetadata  `xml:"report_metadata"`
	Policy   PolicyPublished `xml:"policy_published"`
	Records  []Record        `xml:"record"`
}
