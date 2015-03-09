package wiretap

import "net/http"

type APIServer struct {
}

// NewAPIServer returns a new APIServer
// TODO config dis
func NewAPIServer() *APIServer {
	api := APIServer{}
	return &api
}

func (api *APIServer) ListenAndServe(addr string) error {
	s := &http.Server{
		Addr:    addr,
		Handler: api,
	}
	return s.ListenAndServe()
}

func (api APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(`{"requests":[{"id":1,"url":"www.google.com","returnCode":200}]}`))
}
