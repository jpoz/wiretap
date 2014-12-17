package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jpoz/wiretap"
)

func main() {
	tr := &wiretap.Transport{}
	client := &http.Client{Transport: tr}
	resp, _ := client.Get("https://www.google.com/search?q=jpoz")

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("\n Output: %s", bodyBytes)
}
