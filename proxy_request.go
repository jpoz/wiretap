package wiretap

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

type Proxy struct {
	Client   *http.Client
	Director func(*http.Request)
}

func NewProxy(client *http.Client) Proxy {
	return Proxy{
		Client:   client,
		Director: func(r *http.Request) {},
	}
}

func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	output, _ := ioutil.ReadAll(r.Body)

	reader := bytes.NewReader(output)
	readerCloser := ioutil.NopCloser(reader)

	req, err := http.NewRequest(r.Method, r.URL.String(), readerCloser)
	if err != nil {
		log.Println(err)
		return
	}

	p.Director(req)

	resp, err := p.Client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(bodyBytes)
}
