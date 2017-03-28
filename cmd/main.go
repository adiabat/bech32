package main

import (
	"crypto/rand"
	"fmt"

	"github.com/adiabat/bech32"
)

func main() {

	sr := make([]byte, 22)
	_, _ = rand.Read(sr)

	fmt.Printf("%x\n", sr)

	five := bech32.Bytes8to5(sr)
	fmt.Printf("%x\n", five)

	fivestring, err := bech32.EncodeString(five)
	if err != nil {
		panic(err)
	}
	fmt.Printf("len %d encoded: %s\n", len(fivestring), fivestring)

	fiveback, err := bech32.DecodeString(fivestring)
	if err != nil {
		panic(err)
	}
	fmt.Printf("back to %x\n", fiveback)
	eight, err := bech32.Bytes5to8(five)
	if err != nil {
		panic(err)
	}
	fmt.Printf("back to %x\n", eight)

	hrp := "bc"

	chk := bech32.CreateChecksum(hrp, five)
	// append checksum to end
	fivesum := append(five, chk...)

	fmt.Printf("data+checksum %x\n", fivesum)
	worked := bech32.VerifyChecksum(hrp, fivesum)
	fmt.Printf("checksum %v\n", worked)
	return
}
