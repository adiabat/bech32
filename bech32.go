package bech32

import "fmt"

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

func Mech() string {
	s := fmt.Sprintf("%c", charset[2])
	return s
}

// Bech32PolyModStep operates on 4 bytes at a time...
func Bech32PolyModStep(pre uint32) uint32 {
	gen := []uint32{
		0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3,
	}

	b := uint8(pre >> 25)

	chk := (pre & 0x1fffffff) << 5

	for i, g := range gen {
		chk ^= -(uint32(b>>uint8(i)) & 1) & g
	}
	return chk
}

func Bech32HRPExpand(input string) []byte {
	output := make([]byte, (len(input)*2)+1)

	// first half is the input string shifted down 5 bits.
	// not much is going on there in terms of data / entropy
	for i, c := range input {
		output[i] = uint8(c) >> 5
	}
	// then there's a 0 byte separator
	// don't need to set 0 byte in the middle, as it starts out that way

	// second half is the input string, with the top 3 bits zeroed.
	// most of the data / entropy will live here.
	for i, c := range input {
		output[i+len(input)+1] = uint8(c) & 0x1f
	}
	return output
}

func Bech32CreateChecksum(hrp string, data []byte) []byte {
	values := append(Bech32HRPExpand(hrp), data...)
	// put 6 zero bytes on at the end
	values = append(values, make([]byte, 6)...)
	//	iterate through values applying polyMod

	//checksum := Bech32PolyModStep(values)
	// flip the LSB for good luck
	//	checksum ^= 1
	return nil
}

func Bech32VerifyChecksum(hrp string, data []byte) bool {
	return false
}

func Bech32Encode(hrp string, data []byte) string {
	return "test"
}

func Bech32Decode(hrp string, data []byte) ([]byte, error) {

	return nil, nil
}

//
