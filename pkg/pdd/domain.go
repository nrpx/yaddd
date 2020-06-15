package pdd

import (
	"strconv"

	"golang.org/x/sync/errgroup"
)

const onPage = 10

type Domains struct {
	PddResult `json:",inline"`
	Page      int
	OnPage    int `json:"on_page"`
	Total     int
	Found     int
	Domains   []Domain
}

type Domain struct {
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

func (c *Client) GetDomains() (d []Domain, err error) {
	first, err := c.getDomainsPage(1)
	if err != nil {
		return
	}

	d = append(d, first.Domains...)

	if first.Total <= first.OnPage {
		return
	}

	var errGr errgroup.Group

	pages := (first.Total + first.OnPage - 1) / first.OnPage
	domainsCh := make(chan Domains, pages-1)

	for i := 2; i <= pages; i++ {
		i := i // For `pages > GOMAXPROCS`

		errGr.Go(func() error {
			return c.getDomainsRoutine(domainsCh, i)
		})
	}

	if err = errGr.Wait(); err != nil {
		return nil, err
	}

	close(domainsCh)

	for domains := range domainsCh {
		d = append(d, domains.Domains...)
	}

	return
}

func (c *Client) getDomainsPage(page int) (d Domains, err error) {
	p := Params{
		"page":    strconv.Itoa(page),
		"on_page": strconv.Itoa(onPage),
	}

	req := Request{c, serviceDomain, methodDomains, p}

	d = Domains{}
	if err = req.do(&d); err != nil {
		return
	}

	return
}

func (c *Client) getDomainsRoutine(ch chan Domains, page int) (err error) {
	domains, err := c.getDomainsPage(page)
	if err != nil {
		return
	}

	ch <- domains

	return
}
