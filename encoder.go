package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// interface that implements EncodeRLP(any) []byte
// encodeRLP takes data of any type
// inside it swtiches on type (using the `reflect` pkg)
// there are appropriate handlers for
// single byte vals
// byte arrs and strings
// lists i.e. slices of any type + structs

type RLPEncoder interface {
	EncodeRLP(data interface{}) []byte
}

type Enc struct{}

// encode the DS into bytes slice and get its size
func getSize(data interface{}) int {
	v := reflect.ValueOf(data)
	return int(unsafe.Sizeof(v))
}

func intsize(i uint64) (size int) {
	for size = 1; ; size++ {
		if i >>= 8; i == 0 {
			return size
		}
	}
}

func (e *Enc) EncodeRLP(data interface{}) []byte {
	t := reflect.ValueOf(data)

	fmt.Printf("type is %v\n", t.Kind())

	switch t.Kind() {
	case reflect.Int:
		if t.Int() < 128 {
			return []byte{byte(t.Int())}
		}
		return append(encodeLength(intsize(uint64(t.Int())), 0x80), []byte{byte(t.Int())}...)
	}

	return nil
}

func encodeLength(length int, offset byte) []byte {
	return []byte{byte(length) + offset}
}
