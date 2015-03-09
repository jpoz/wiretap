package main

import (
	"fmt"
	"net/http"

	"github.com/jpoz/wiretap"
	"github.com/jpoz/wiretap/disk"
)

// Edit
func director(r *http.Request) {}

func main() {
	// Create wiretap Transport
	tr := &wiretap.Transport{
		Storage: disk.Storage{"./cache"},
	}

	// HTTP client with wiretap Transport
	client := &http.Client{Transport: tr}

	// Proxy
	proxy := wiretap.HttpProxy{
		Client:   client,
		Director: director,
	}

	s := &http.Server{
		Addr:    ":9080",
		Handler: proxy,
	}

	go func() {
		println("serving")
		err := s.ListenAndServe()
		fmt.Println(err)
	}()

	socks, err := wiretap.NewSocksProxy()
	if err != nil {
		fmt.Printf("err %+v\n", err)
	}

	println("Sock Started")
	if err := socks.ListenAndServe("tcp", "127.0.0.1:8000"); err != nil {
		panic(err)
	}
	fmt.Println(err)
}
