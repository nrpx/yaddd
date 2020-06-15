package pdd

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// Параметры запроса в виде пары "ключ-значение".
type Params map[string]string

// Запрос к API.
type Request struct {
	client  *Client
	service Service
	method  Method
	params  Params
}

type Response interface {
	Result() (string, error)
}

// Функция возвращает готовый URL для указанных раздела и метода.
func getURLFor(s Service, m Method, p Params) (u url.URL, err error) {
	if s == 0 {
		return u, serviceUnknownErr
	}

	if m == 0 {
		return u, methodUnknownErr
	}

	q := url.Values{}

	for key, val := range p {
		q.Set(key, val)
	}

	u = url.URL{
		Scheme:   pddScheme,
		Host:     pddHost,
		Path:     path.Join(pddAPIPath, s.String(), m.String()),
		RawQuery: q.Encode(),
	}

	return
}

func (r Request) do(rr Response) (err error) {
	u, err := getURLFor(r.service, r.method, r.params)
	if err != nil {
		return
	}

	req, err := http.NewRequest(r.method.HTTPMethod(), u.String(), nil)
	if err != nil {
		return
	}

	req.Header.Add("PddToken", r.client.pddToken)

	resp, err := r.client.Do(req)
	if err != nil {
		return
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	decoder := json.NewDecoder(resp.Body)
	if err != nil {
		return
	}

	if err = decoder.Decode(rr); err != nil {
		return
	}

	if _, err = rr.Result(); err != nil {
		return
	}

	return
}
