package wiretap

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Storage implments the storage of a wiretap.Session
type Storage interface {
	Save(s Session)
}

// Session is an entire request/response cycle
type Session struct {
	Started     time.Time
	Request     *http.Request
	RequestBody []byte

	Response     *http.Response
	ResponseBody []byte
	Completed    time.Time
}

// Transport is an http.RoundTripper and can be used inplace of the
// http.Transport
type Transport struct {
	Storage Storage
	http.Transport
}

func tap(input io.ReadCloser) ([]byte, io.ReadCloser) {
	output, _ := ioutil.ReadAll(input)

	reader := bytes.NewReader(output)
	readerCloser := ioutil.NopCloser(reader)

	return output, readerCloser
}

// RoundTrip actually use the http.Transport.RoundTrip function but
// Reads the Request body before and the Response body after
func (t Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	session := Session{Request: r, Started: time.Now()}
	if r.Body != nil {
		session.RequestBody, r.Body = tap(r.Body)
	}

	var err error
	session.Response, err = t.Transport.RoundTrip(r)
	if err != nil {
		return session.Response, err
	}
	session.ResponseBody, session.Response.Body = tap(session.Response.Body)

	if t.Storage != nil {
		go t.Storage.Save(session)
	}

	session.Completed = time.Now()
	return session.Response, err
}
