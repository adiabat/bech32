package bech32

import (
	"bytes"
	"math/rand"
	"testing"
)

// TestRandom makes random signatures and compresses / decompresses them
func TestRandomBytesBackAndForth(t *testing.T) {

	seed := make([]byte, 20)
	_, _ = rand.Read(seed)

	reduceBits := Bytes8to5(seed)

	fivestring, err := EncodeString(reduceBits)
	if err != nil {

	}

	fiveback, err := DecodeString(fivestring)
	if err != nil {
		t.Fatal(err)
	}

	eight, err := Bytes5to8(reduceBits)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(fiveback, reduceBits) {
		t.Fatalf("reencode mismatch %x != %x", fiveback, reduceBits)
	}

	if !bytes.Equal(seed, eight) {
		t.Fatalf("reencode mismatch %x != %x", fiveback, reduceBits)
	}

}
