package main

import (
	"fmt"
	"net/http"

	"github.com/jpoz/wiretap"
	"github.com/jpoz/wiretap/disk"
)

type BasicHandler struct{}

func (h BasicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	w.Write([]byte("OK"))
}

const Scheme = "http"
const Host = "github.com"

// Edit
func director(r *http.Request) {
	r.URL.Scheme = Scheme
	r.URL.Host = Host
}

func main() {
	// Create wiretap Transport
	tr := &wiretap.Transport{
		Storage: disk.Storage{"./cache"},
	}

	// HTTP client with wiretap Transport
	client := &http.Client{Transport: tr}

	// Proxy
	proxy := wiretap.Proxy{
		Client:   client,
		Director: director,
	}

	s := &http.Server{
		Addr:    ":9898",
		Handler: proxy,
	}

	go func() {
		println("serving")
		err := s.ListenAndServe()
		fmt.Println(err)
	}()

	socks := wiretap.Socks{
		Server: s,
	}

	println("Sock Started")
	err := socks.Start("localhost:8888")
	fmt.Println(err)
}
