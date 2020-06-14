package pdd

import (
	"strconv"

	"golang.org/x/sync/errgroup"
)

const onPage = 10

type PddDomains struct {
	PddResult `json:",inline"`
	Page      int
	OnPage    int `json:"on_page"`
	Total     int
	Found     int
	Domains   []PddDomain
}

type PddDomain struct {
	Name           string
	Status         string
	Aliases        []string
	LogoEnabled    PddBool `json:"logo_enabled"`
	LogoURL        string  `json:"log_url"`
	NSDelegated    string  `json:"nsdelegated"`
	MasterAdmin    PddBool `json:"master_admin"`
	DKIMReady      PddBool `json:"dkim-ready"`
	EmailsMaxCount int     `json:"emails-max-count"`
	EmailsCount    int     `json:"emails-count"`
	NoDKIM         PddBool `json:"nodkim"`
}

func (c *Client) GetDomains() (d []PddDomain, err error) {
	first, err := c.getDomainsPage(1)
	if err != nil {
		return
	}

	if first.Total > first.OnPage {
		pages := first.Total / first.OnPage

		if first.Total%first.OnPage > 0 {
			pages += 1
		}

		var errGr errgroup.Group
		domainsCh := make(chan PddDomains, pages-1)

		for i := 2; i <= int(pages); i++ {
			errGr.Go(func() error {
				return c.getDomainsRoutine(domainsCh, i)
			})
		}

		if err = errGr.Wait(); err != nil {
			return nil, err
		}

		for domains := range domainsCh {
			d = append(d, domains.Domains...)
		}
	}

	d = append(first.Domains, d...)

	return
}

func (c *Client) getDomainsPage(page int) (d PddDomains, err error) {
	p := make(Params)
	p["page"] = strconv.Itoa(page)
	p["on_page"] = strconv.Itoa(onPage)

	req := Request{c, serviceDomain, methodDomains, p}

	d = PddDomains{}
	if err = req.do(&d); err != nil {
		return
	}

	return
}

func (c *Client) getDomainsRoutine(ch chan PddDomains, page int) (err error) {
	domains, err := c.getDomainsPage(page)
	if err != nil {
		return
	}

	ch <- domains

	return
}
