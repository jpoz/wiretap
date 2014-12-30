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
		Storage: disk.Storage{"./cache"},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("http://jsonip.com/")

	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("bodyBytes %+s\n", bodyBytes)

	time.Sleep(100 * time.Millisecond)
}
