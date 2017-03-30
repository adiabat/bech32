package bech32

import (
	"bytes"
	"math/rand"
	"testing"
)

type validTestAddress struct {
	address string
	data    []byte
}

var (
	validChecksum = []string{
		"A12UEL5L",
		"an83characterlonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1tt5tgs",
		"abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw",
		"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j",
		"split1checkupstagehandshakeupstreamerranterredcaperred2y9e3w",
	}

	invalidAddress = []string{
		"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kg3g4ty",
		"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t5",
		"BC13W508D6QEJXTDG4Y5R3ZARVARY0C5XW7KN40WF2",
		"bc1rw5uspcuh",
		"bc10w508d6qejxtdg4y5r3zarvary0c5xw7kw508d6qejxtdg4y5r3zarvary0c5xw7kw5rljs90",
		"BC1QR508D6QEJXTDG4Y5R3ZARVARYV98GJ9P",
		"tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sL5k7",
		"tb1pw508d6qejxtdg4y5r3zarqfsj6c3",
		"tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3pjxtptv",
	}
)

func TestRandomEncodeDecode(t *testing.T) {
	data := make([]byte, 20)
	_, _ = rand.Read(data)

	hrp := "test"

	adr := Encode(hrp, data)

	hrp2, data2, err := Decode(adr)
	if err != nil {
		t.Fatal(err)
	}
	if hrp2 != hrp {
		t.Fatalf("hrp mismatch %s, %s", hrp, hrp2)
	}
	if !bytes.Equal(data, data2) {
		t.Fatalf("data mismatch %x, %x", data, data2)
	}

}

func TestHardCoded(t *testing.T) {

	for _, adr := range validChecksum {
		_, _, err := Decode(adr)
		if err != nil {
			t.Fatalf("address %s invalid:%s", adr, err)
		}
	}

	for _, adr := range invalidAddress {
		_, _, err := Decode(adr)
		if err == nil {
			t.Fatalf("address %s should fail but didn't", adr)
		}
	}

}
