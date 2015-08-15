package main

import (
	"log"

	"github.com/jpoz/wiretap"
)

func main() {
	wiretap, _ := wiretap.NewWiretap()
	log.Fatal(wiretap.ListenAndServe())
}
