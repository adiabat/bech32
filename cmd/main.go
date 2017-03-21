package main

import (
	"fmt"

	"github.com/adiabat/bech32"
)

func main() {
	s := bech32.Mech()
	fmt.Printf("%s\n", s)
	return
}
