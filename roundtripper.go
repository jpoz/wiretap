package wiretap

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type Session struct {
	Request     *http.Request
	RequestBody []byte

	Response     *http.Response
	ResponseBody []byte
}

type Transport struct {
	http.Transport
}

func Tap(input io.ReadCloser) ([]byte, io.ReadCloser) {
	output, _ := ioutil.ReadAll(input)

	reader := bytes.NewReader(output)
	readerCloser := ioutil.NopCloser(reader)

	return output, readerCloser
}

func (t Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	session := Session{Request: r}
	if r.Body != nil {
		session.RequestBody, r.Body = Tap(r.Body)
	}

	var err error
	session.Response, err = t.Transport.RoundTrip(r)
	session.ResponseBody, session.Response.Body = Tap(session.Response.Body)

	return session.Response, err
}
