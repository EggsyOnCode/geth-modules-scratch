package main

import (
	"fmt"
	"reflect"
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
	h := []byte{0x68, 0x89}
	fmt.Printf("elem type is %v\n", reflect.TypeOf(h).Elem())
	exp := []byte{0x82, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
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
