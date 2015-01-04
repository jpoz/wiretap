package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jpoz/wiretap"
	"github.com/jpoz/wiretap/disk"
)

const Host = "http://github.com/"

func main() {
	// Create wiretap Transport
	tr := &wiretap.Transport{
		Storage: disk.Storage{"./cache"},
	}

	// HTTP client with wiretap Transport
	client := &http.Client{Transport: tr}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest("GET", Host+r.URL.Path, nil)
		resp, err := client.Do(req)

		if err != nil {
			log.Println(err)
			w.WriteHeader(resp.StatusCode)
			return
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		w.Write(bodyBytes)
		log.Println("Wrote to ./cache")
	})

	fmt.Println("localhost:8000 -> http://example.com/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
