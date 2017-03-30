package main

import (
	"crypto/rand"
	"fmt"

	"github.com/adiabat/bech32"
)

func main() {
	testAdr := "split1checkupstagehandshakeupstreamerranterredcaperred2y9e3w"

	hrp, decodedData, err := bech32.Decode(testAdr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("decoded data: %x\n", decodedData)

	//	testBytes := make([]byte, 32)
	//	for i, _ := range testBytes {
	//		testBytes[i] = byte(i)
	//	}

	//	fmt.Printf("testBytes: %x\n", testBytes)

	adr := bech32.Encode(hrp, decodedData)

	//	fmt.Printf("testadr: %s\n", testAdr)
	fmt.Printf("start  : %s\n", testAdr)
	fmt.Printf("encoded: %s\n", adr)

}

func main1() {

	sr := make([]byte, 22)
	_, _ = rand.Read(sr)

	fmt.Printf("%x\n", sr)

	five := bech32.Bytes8to5(sr)
	fmt.Printf("%x\n", five)

	return
}

/*
	testBytes := make([]byte, 32)
	for i, _ := range testBytes {
		testBytes[i] = byte(i)
	}
*/
