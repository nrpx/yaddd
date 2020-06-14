package core

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"yaddd/internal/config"
	"yaddd/internal/core/pdd"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var (
	domainNotFoundErr    = errors.New("Domain not found")
	aRecordNotFoundErr   = errors.New("A-record not found")
	moreThanOneRecordErr = errors.New("Found more than one A-record")
)

type dynDNS struct {
	pddClient *pdd.Client
	conf      *config.Config
}

func StartService(conf *config.Config) (err error) {
	pddClient, err := pdd.NewClient(conf.PddToken)
	if err != nil {
		logrus.WithError(err).Fatal("Create PDD client")
	}

	d := &dynDNS{pddClient, conf}

	if err = d.checkDomain(); err != nil {
		logrus.WithError(err).Fatal("Check domain")
	}

	c := cron.New()

	id, err := c.AddFunc(conf.Cron, d.updateIP)
	if err != nil {
		logrus.WithError(err).Fatal("Add cron job")
	}

	logrus.WithField("cron-id", id).Info("Cron job added")

	go c.Start()

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	return
}

func (d *dynDNS) checkDomain() (err error) {
	domains, err := d.pddClient.GetDomains()
	if err != nil {
		return
	}

	var ok bool
	for _, domain := range domains {
		if domain.Name == d.conf.DomainName {
			ok = true
			break
		}
	}

	if !ok {
		return domainNotFoundErr
	}

	return
}

func (d *dynDNS) getARecord() (r pdd.DNSRecordStruct, err error) {
	records, err := d.pddClient.GetDNSRecords(d.conf.DomainName)
	if err != nil {
		return
	}

	aRecords := records.FilterByType(pdd.TypeA)

	switch d.conf.CurrentIP {
	case "":
		if len(aRecords) > 1 {
			return r, moreThanOneRecordErr
		} else {
			r = aRecords[0]
		}
	default:
		for _, a := range aRecords {
			if a.Content == d.conf.CurrentIP {
				r = a
				break
			}
		}
	}

	if r.Domain == "" {
		return r, aRecordNotFoundErr
	}

	return
}

func (d *dynDNS) updateIP() {
	ip, err := GetExternalIP()
	if err != nil {
		logrus.WithError(err).Error("Get external IP")

		return
	}

	params := pdd.DNSRecordParams{
		"domain":  d.conf.DomainName,
		"content": ip,
	}

	record, err := d.getARecord()

	var resp pdd.DNSRecord

	switch {
	case record.Content == ip:
		logrus.WithField("record-id", record.RecordID).
			WithField("record-ip", record.Content).
			Info("IP is relevant")

		return
	case errors.Is(err, aRecordNotFoundErr):
		resp, err = d.pddClient.AddDNSRecord(pdd.TypeA, params)
	case err != nil:
		logrus.WithError(err).Error("Get A-record")

		return
	default:
		resp, err = d.pddClient.EditDNSRecord(record.RecordID, params)
	}

	if err != nil {
		logrus.WithError(err).Error("Create/update A-record")

		return
	}

	d.conf.CurrentIP = resp.Record.Content

	logrus.WithField("record-id", resp.Record.RecordID).
		WithField("ip-address", ip).
		Debug("Record created/updated")
}
