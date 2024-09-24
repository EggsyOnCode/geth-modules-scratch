package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSize(t *testing.T) {
	str := 124
	expSize := 1
	size := intsize(uint64(str))
	fmt.Printf("size is %v", size)
	assert.Equal(t, size, expSize)
}

func TestEncodeForUin8(t *testing.T) {
	enc := Enc{}
	data := 178
	exp := []byte{0x81, 0xb2}
	res := enc.EncodeRLP(data)
	fmt.Printf("res is %v\n", res)
	assert.Equal(t, res, exp)
}

func TestEncodeForString(t *testing.T) {
	enc := Enc{}
	data := "hello"
	exp := []byte{0x85, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
	res := enc.EncodeRLP(data)
	fmt.Printf("res is %v\n", res)
	assert.Equal(t, res, exp)
}

func TestEncodeByteSlice(t *testing.T) {
	enc := Enc{}
	data := []byte{0x01, 0x02, 0x03}
	exp := []byte{0x83, 0x01, 0x02, 0x03}
	res := enc.EncodeRLP(data)
	fmt.Printf("res is %v\n", res)
	assert.Equal(t, res, exp)
}

func TestReflectDT(t *testing.T) {
	enc := Enc{}
	data := []string{"hello", "world"}
	exp := []byte{0xcc, 0x85, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x85, 0x77, 0x6f, 0x72, 0x6c, 0x64}
	res := enc.EncodeRLP(data)
	fmt.Printf("res is %v\n", res)
	assert.Equal(t, res, exp)

	// the prefix exceeds 0xf8 - 0xff limit if the data set is large
}
