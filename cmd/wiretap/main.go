package main

import (
	"log"

	"github.com/jpoz/wiretap"
)

func main() {
	wiretap := wiretap.NewWiretap()
	log.Fatal(wiretap.ListenAndServe())
}
