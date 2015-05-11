package wiretap

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/jpoz/env/decoder"
)

type APIServer struct {
	Addr string `expand:"${APIHOST}:${APIPORT}"`
}

// NewAPIServer returns a new APIServer
func NewAPIServer() *APIServer {
	api := APIServer{}
	decoder.Decode(&api)
	return &api
}

func (api *APIServer) ListenAndServe() error {
	s := &http.Server{
		Addr:    api.Addr,
		Handler: api,
	}

	log.Infof("Starting API server on %s", api.Addr)
	return s.ListenAndServe()
}

func (api APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(`{"requests":[{"id":1,"url":"www.google.com","returnCode":200}]}`))
}
