package bech32

import "fmt"

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

func Mech() string {
	s := fmt.Sprintf("%c", charset[2])
	return s
}

// Bytes8to5 extends a byte slice into a longer, padded byte slice of 5-bit elements
// where the high 3 bits are all 0.
// I guess you never need to pad...?
func Bytes8to5(input []byte) []byte {

	// round up divistion for output length
	outputLength := (((len(input)) * 8) + 4) / 5

	dest := make([]byte, outputLength)

	outputOffset := 0

	// Continue until input has been consumed
	for len(input) > 0 {

		// 5 input bytes, map to 8 output bytes
		// got switches from base32 stdlib
		switch len(input) {
		default:
			dest[outputOffset+7] = input[4] & 0x1f
			dest[outputOffset+6] = input[4] >> 5
			fallthrough
		case 4:
			dest[outputOffset+6] |= (input[3] << 3) & 0x1f
			dest[outputOffset+5] = (input[3] >> 2) & 0x1f
			dest[outputOffset+4] = input[3] >> 7
			fallthrough
		case 3:
			dest[outputOffset+4] |= (input[2] << 1) & 0x1f
			dest[outputOffset+3] = (input[2] >> 4) & 0x1f
			fallthrough
		case 2:
			dest[outputOffset+3] |= (input[1] << 4) & 0x1f
			dest[outputOffset+2] = (input[1] >> 1) & 0x1f
			dest[outputOffset+1] = (input[1] >> 6) & 0x1f
			fallthrough
		case 1:
			dest[outputOffset+1] |= (input[0] << 2) & 0x1f
			dest[outputOffset+0] = input[0] >> 3
		}

		if len(input) < 5 {
			break
		}

		// pop off first 5 bytes of input slice as we're done with them
		input = input[5:]
		// our output slice has advanced 8 bytes for the next round
		outputOffset += 8
	}

	return dest
}

func Bytes5to8(input []byte) ([]byte, error) {

	// first check that high 3 bits for all bytes are 0
	for i, b := range input {
		if b&0xe0 != 0 {
			return nil, fmt.Errorf("Invalid byte at position %d: %x", i, b)
		}
	}

	outputOffset := 0

	// round up and divide to get output length
	dest := make([]byte, (((len(input) * 5) + 7) / 8))

	for len(input) > 0 {

		switch len(input) {
		default:
			dest[outputOffset+4] = input[6]<<5 | input[7]
			fallthrough
		case 7:
			dest[outputOffset+3] = input[4]<<7 | input[5]<<2 | input[6]>>3
			fallthrough
		case 5:
			dest[outputOffset+2] = input[3]<<4 | input[4]>>1
			fallthrough
		case 4:
			dest[outputOffset+1] = input[1]<<6 | input[2]<<1 | input[3]>>4
			fallthrough
		case 2:
			dest[outputOffset] = input[0]<<3 | input[1]>>2
		}

		// if there is less than 8 characters left in the input string, we're done
		if len(input) < 8 {
			break
		}

		// pop off first 8 characters of the 32-bit encoded string
		input = input[8:]
		// advance output position by 5 bytes
		outputOffset += 5
	}

	return dest, nil
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
