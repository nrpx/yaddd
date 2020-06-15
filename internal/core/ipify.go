package core

import (
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

const ipifyURL = "https://api.ipify.org"

func GetExternalIP() (ip string, err error) {
	res, err := http.Get(ipifyURL)
	if err != nil {
		logrus.WithError(err).Error("Get external IP")

		return
	}

	defer res.Body.Close()

	rawIP, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.WithError(err).Error("Read received IP")

		return
	}

	return string(rawIP), err
}
