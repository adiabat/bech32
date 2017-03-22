package main

import (
	"fmt"

	"github.com/adiabat/bech32"
)

func main() {

	sr := []byte("XokokOKOKOxxzzzzwwwzzzzzzzzX--@@@!!}{--")
	five := bech32.Bytes8to5(sr)
	fmt.Printf("%x\n", five)
	eight, err := bech32.Bytes5to8(five)
	if err != nil {
		panic(err)
	}
	fmt.Printf("back to %x which is:\n%s\n", eight, eight)

	return
}
