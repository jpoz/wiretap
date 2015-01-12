package main

import (
	"fmt"
	"net/http"

	"github.com/jpoz/wiretap"
	"github.com/jpoz/wiretap/disk"
)

func main() {
	// Works like a normal http.Client but writes Request/Response to "./cache"
	tr := &wiretap.Transport{
		Storage: disk.Storage{"./cache"},
	}
	client := &http.Client{Transport: tr}

	resp, _ := client.Get("http://jsonip.com/")

	fmt.Printf("Returned %+s\n", resp.Status)

	fmt.Println("Wrote to ./cache/jsonpi.com/GET/{TIMESTAMP}/response.txt")
}
