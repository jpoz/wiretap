package main

import (
	"log"

	"github.com/jpoz/wiretap"
)

func main() {
	api := wiretap.NewAPIServer()
	log.Println("Running on :8000")
	log.Fatal(api.ListenAndServe(":8000"))
}
