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
