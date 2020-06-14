package pdd

import (
	"errors"
	"net/http"
)

var pddTokenErr = errors.New("Incorrect token")

// Client is an http.Client for Yandex Connect API
type Client struct {
	*http.Client
	pddToken string
}

func NewClient(token string) (c *Client, err error) {
	if token == "" {
		return nil, pddTokenErr
	}

	c = &Client{&http.Client{}, token}

	return
}
