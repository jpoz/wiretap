package main

import (
	"fmt"

	"github.com/jpoz/wiretap/socks"
)

func main() {
	proxy, _ := socks.NewProxy("tcp", ":8888")
	fmt.Println(proxy.ListenAndServe())
}
