package bech32

import (
	"fmt"
	"strings"
)

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var inverseCharset = [256]int8{
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	15, -1, 10, 17, 21, 20, 26, 30, 7, 5, -1, -1, -1, -1, -1, -1,
	-1, 29, -1, 24, 13, 25, 9, 8, 23, -1, 18, 22, 31, 27, 19, -1,
	1, 0, 3, 16, 11, 28, 12, 14, 6, 4, 2, -1, -1, -1, -1, -1,
	-1, 29, -1, 24, 13, 25, 9, 8, 23, -1, 18, 22, 31, 27, 19, -1,
	1, 0, 3, 16, 11, 28, 12, 14, 6, 4, 2, -1, -1, -1, -1, -1}

// Bytes8to5 extends a byte slice into a longer, padded byte slice of 5-bit elements
// where the high 3 bits are all 0.
// I guess you never need to pad...?
func Bytes8to5(input []byte) []byte {

	// round up divistion for output length
	outputLength := (((len(input)) * 8) + 4) / 5

	dest := make([]byte, outputLength)

	outputOffset := 0

	//	fmt.Printf("8 to 5 input len %d output len %d\n", len(input), len(dest))
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

	//	fmt.Printf("5 to 8 input len %d output len %d\n", len(input), len(dest))

	for len(input) > 0 {

		switch len(input) {
		default:
			dest[outputOffset+4] = input[7]
			fallthrough
		case 7:
			dest[outputOffset+4] |= input[6] << 5
			dest[outputOffset+3] = input[6] >> 3
			fallthrough
		case 6:
			dest[outputOffset+3] |= input[5] << 2
			fallthrough
		case 5:
			dest[outputOffset+3] |= input[4] << 7
			dest[outputOffset+2] = input[4] >> 1
			fallthrough
		case 4:
			dest[outputOffset+2] |= input[3] << 4
			dest[outputOffset+1] = input[3] >> 4
			fallthrough
		case 3:
			dest[outputOffset+1] |= input[2] << 1
			fallthrough
		case 2:
			dest[outputOffset+1] |= input[1] << 6
			dest[outputOffset] = input[1] >> 2
			fallthrough
		case 1:
			dest[outputOffset] |= input[0] << 3
		}

		// if there are fewer than 8 characters left in the input string, we're done
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

// EncodeString swaps 5-bit bytes with a string of the corresponding letters
func BytesToString(input []byte) (string, error) {
	var s string
	for i, c := range input {
		if c&0xe0 != 0 {
			return "", fmt.Errorf("high bits set at position %d: %x", i, c)
		}
		s += string(charset[c])
	}
	return s, nil
}

func StringToBytes(input string) ([]byte, error) {
	b := make([]byte, len(input))
	for i, c := range input {
		if inverseCharset[c] == -1 {
			return nil, fmt.Errorf("contains invalid character %s", string(c))
		}
		b[i] = byte(inverseCharset[c])
	}
	return b, nil
}

// PolyMod takes a slice and give the polymod *uint32*.  I think
func PolyMod(values []byte) uint32 {
	gen := []uint32{
		0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3,
	}

	chk := uint32(1)

	for _, v := range values {
		top := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ uint32(v)
		for i, g := range gen {
			if (top>>uint8(i))&1 == 1 {
				chk ^= g
			}
		}
	}

	return chk
}

// HRPExpand turns the human redable part into 5bit-bytes for later processing
func HRPExpand(input string) []byte {
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

// create checksum makes a 6-shortbyte checksum from the HRP and data parts
func CreateChecksum(hrp string, data []byte) []byte {
	values := append(HRPExpand(hrp), data...)
	// put 6 zero bytes on at the end
	values = append(values, make([]byte, 6)...)
	//get checksum for whole slice

	checksum := PolyMod(values) ^ 1

	for i := 0; i < 6; i++ {
		// note that this is NOT the same as converting 8 to 5
		// this is it's own expansion to 6 bytes from 4, chopping
		// off the MSBs.
		values[(len(values)-6)+i] = byte(checksum>>(5*(5-uint32(i)))) & 0x1f
	}

	return values[len(values)-6:]
}

func VerifyChecksum(hrp string, data []byte) bool {
	values := append(HRPExpand(hrp), data...)
	checksum := PolyMod(values)
	//	fmt.Printf("checksum %x\n", checksum)
	return checksum == 1
}

func Encode(hrp string, data []byte) string {
	fiveData := Bytes8to5(data)
	return EncodeSquashed(hrp, fiveData)
}

func EncodeSquashed(hrp string, data []byte) string {
	combined := append(data, CreateChecksum(hrp, data)...)

	// ignore error as we just five'd it
	dataString, err := BytesToString(combined)
	if err != nil {
		panic(err)
	}

	return hrp + "1" + dataString
}

func Decode(adr string) (string, []byte, error) {
	hrp, squashedData, err := DecodeSquashed(adr)
	if err != nil {
		return "", nil, err
	}
	data, err := Bytes5to8(squashedData)
	if err != nil {
		return "", nil, err
	}
	return hrp, data, nil
}

func DecodeSquashed(adr string) (string, []byte, error) {

	lowAdr := strings.ToLower(adr)
	highAdr := strings.ToUpper(adr)

	if adr != lowAdr && adr != highAdr {
		return "", nil, fmt.Errorf("mixed case address")
	}

	adr = lowAdr

	// find the last "1" and split there
	splitLoc := strings.LastIndex(adr, "1")
	if splitLoc == -1 {
		return "", nil, fmt.Errorf("1 separator not present in address")
	}
	hrp := adr[0:splitLoc]

	data, err := StringToBytes(adr[splitLoc+1:])
	if err != nil {
		return "", nil, err
	}

	sumOK := VerifyChecksum(hrp, data)
	if !sumOK {
		return "", nil, fmt.Errorf("Checksum invalid")
	}
	data = data[:len(data)-6]

	return hrp, data, nil
}

func SegWitAddressEncode(hrp string, data []byte) (string, error) {
	//	combined := append(data, CreateChecksum(hrp, data))
	return "", nil
}

func SegWitAddressDecode(adr string) ([]byte, error) {
	hrp, squashedData, err := DecodeSquashed(adr)
	if err != nil {
		return nil, err
	}
	// the segwit version byte is directly put into a 5bit squashed byte
	// since it maxes out at 16, wasting ~1 byte instead of 4.

	version := squashedData[0]
	data, err := Bytes5to8(squashedData[1:])
	if err != nil {
		return nil, err
	}
	if hrp != "bc" && hrp != "tb" {
		return nil, fmt.Errorf("prefix %s is not bitcoin or testnet", hrp)
	}
	if len(data) < 2 || len(data) > 40 {
		return nil, fmt.Errorf("Data length %d out of bounds", len(data))
	}

	if version > 16 {
		return nil, fmt.Errorf("Invalid witness program version %d", data[0])
	}
	if version == 0 && len(data) != 20 && len(data) != 32 {
		return nil, fmt.Errorf("expect 20 or 32 byte v0 witprog, got %d", len(data))
	}
	return append([]byte{version}, data...), nil
}

//
