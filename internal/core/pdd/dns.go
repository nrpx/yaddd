package pdd

import (
	"errors"
	"strconv"
)

type DNSRecordType string

const (
	TypeA     DNSRecordType = "A"
	TypeAAAA  DNSRecordType = "AAAA"
	TypeCNAME DNSRecordType = "CNAME"
	TypeMX    DNSRecordType = "MX"
	TypeNS    DNSRecordType = "NS"
	TypeSOA   DNSRecordType = "SOA"
	TypeSRV   DNSRecordType = "SRV"
	TypeTXT   DNSRecordType = "TXT"
)

var domainEmptyErr = errors.New("Domain is required")

type DNSRecordStruct struct {
	RecordID  int `json:"record_id"`
	Type      string
	Domain    string
	FQDN      string
	TTL       PddInt
	Subdomain string
	Content   string
	Priority  PddInt
}

type DNSRecords struct {
	PddResult `json:",inline"`
	Domain    string
	Records   []DNSRecordStruct
}

type DNSRecord struct {
	PddResult `json:",inline"`
	Domain    string
	Record    DNSRecordStruct
}

// Карта параметров для отправки с запросом к API.
type DNSRecordParams map[string]string

func (p DNSRecordParams) Validate() error {
	if p["domain"] == "" {
		return domainEmptyErr
	}

	return nil
}

func (c *Client) GetDNSRecords(domain string) (r DNSRecords, err error) {
	if domain == "" {
		return r, domainEmptyErr
	}

	p := make(Params)
	p["domain"] = domain

	req := Request{c, serviceDNS, methodList, p}

	r = DNSRecords{}
	if err = req.do(&r); err != nil {
		return
	}

	return
}

func (c *Client) AddDNSRecord(t DNSRecordType, d DNSRecordParams) (r DNSRecord, err error) {
	if err = d.Validate(); err != nil {
		return
	}

	d["type"] = string(t)

	req := Request{c, serviceDNS, methodAdd, Params(d)}

	r = DNSRecord{}
	if err = req.do(&r); err != nil {
		return
	}

	return
}

func (c *Client) EditDNSRecord(id int, d DNSRecordParams) (r DNSRecord, err error) {
	if err = d.Validate(); err != nil {
		return
	}

	d["record_id"] = strconv.Itoa(id)

	req := Request{c, serviceDNS, methodEdit, Params(d)}

	r = DNSRecord{}
	if err = req.do(&r); err != nil {
		return
	}

	return
}

func (d *DNSRecords) FilterByType(t DNSRecordType) (rs []DNSRecordStruct) {
	for _, record := range d.Records {
		if DNSRecordType(record.Type) == t {
			rs = append(rs, record)
		}
	}

	return
}
