package pdd

import (
	"encoding/json"
	"errors"
	"strconv"
)

const (
	pddHost    = "pddimp.yandex.ru"
	pddScheme  = "https"
	pddAPIPath = "api2/admin"
)

var (
	serviceUnknownErr = errors.New("Service unknown")
	methodUnknownErr  = errors.New("Method unknown")
)

type PddBool struct {
	Val  bool
	Text string
}

func (b *PddBool) UnmarshalJSON(data []byte) (err error) {
	b.Text = string(data)
	b.Val = b.Text == "yes"

	return
}

type PddInt struct {
	Val  int
	Text string
}

func (i *PddInt) UnmarshalJSON(data []byte) (err error) {
	if err = json.Unmarshal(data, &i.Text); err != nil {
		if err = json.Unmarshal(data, &i.Val); err != nil {
			return
		}
	}

	if i.Text != "" {
		if i.Val, err = strconv.Atoi(i.Text); err != nil {
			return
		}
	} else {
		i.Text = strconv.Itoa(i.Val)
	}

	return
}

type PddResult struct {
	Success string
	Error   string
}

func (d *PddResult) Result() (string, error) {
	switch {
	case d.Error != "":
		return d.Success, errors.New(d.Error)
	default:
		return d.Success, nil
	}
}
