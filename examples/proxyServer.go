package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jpoz/wiretap"
	"github.com/jpoz/wiretap/disk"
)

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

	http.HandleFunc("/", proxy.ServeHTTP)

	fmt.Println("localhost:8000 -> http://github.com/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
