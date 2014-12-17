package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jpoz/wiretap"
	"github.com/jpoz/wiretap/disk"
)

func main() {
	tr := &wiretap.Transport{
		Storage: disk.Storage{},
	}
	client := &http.Client{Transport: tr}
	resp, _ := client.Get("https://www.google.com/search?q=jpoz")

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("\n Output: %s", len(bodyBytes))

	time.Sleep(100 * time.Millisecond)
}
