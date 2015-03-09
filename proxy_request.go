package wiretap

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpProxy struct {
	Client   *http.Client
	Director func(*http.Request)
}

func NewHttpProxy(client *http.Client) HttpProxy {
	return HttpProxy{
		Client:   client,
		Director: func(r *http.Request) {},
	}
}

func (p HttpProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	output, _ := ioutil.ReadAll(r.Body)

	reader := bytes.NewReader(output)
	readerCloser := ioutil.NopCloser(reader)

	fmt.Printf("r %+v\n", r)

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
