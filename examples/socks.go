package main

import (
	"fmt"

	"github.com/jpoz/wiretap"
)

func main() {
	socks := wiretap.Socks{}

	err := socks.Start("localhost:8888")
	fmt.Println(err)
}
