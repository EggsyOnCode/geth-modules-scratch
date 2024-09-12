package main

import (
	"bytes"
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

	kind := t.Kind()

	switch {
	case kind == reflect.Int:
		if t.Int() < 128 {
			return []byte{byte(t.Int())}
		}
		return append(encodeLength(intsize(uint64(t.Int())), 0x80), []byte{byte(t.Int())}...)
	case kind == reflect.String:
		// strings
		return encodeStringSlice(data)
	case kind == reflect.Slice && isByte(t.Type().Elem()):
		// byte slice
		return encodeByteSlice(data)
	case kind == reflect.Slice:
		// slice of interface{} but not byte slice
		// algo has to be recursive
		return encodeCustomSlice(e, data)
	default:
		fmt.Print("rlp: unsupported type")
	}

	return nil
}

func encodeLength(length int, offset byte) []byte {
	return []byte{byte(length) + offset}
}

func encodeCustomSlice(e *Enc, data interface{}) []byte {
	var buffer bytes.Buffer
	slice := reflect.ValueOf(data)

	for i := 0; i < slice.Len(); i++ {
		// Base condition: If elem is not a slice, encode it
		elem := slice.Index(i).Interface()
		fmt.Printf("elem is %T\n", elem)
		if slice.Index(i).Kind() != reflect.Slice {
			buffer.Write(e.EncodeRLP(elem)) // Append the RLP-encoded element to the buffer
		} else {
			// Recursively encode the slice
			buffer.Write(encodeCustomSlice(e, elem))
		}
	}

	fmt.Printf("len of the data is %v\n", len(buffer.Bytes()))

	// Encode the length of the buffer prefixed with the length of the encoded list
	return append(encodeLength(len(buffer.Bytes()), 0xc0), buffer.Bytes()...)
}

func encodeStringSlice(data interface{}) []byte {
	t := reflect.ValueOf(data)
	str := t.String()
	length := len(str)

	if length < 55 {
		// For lengths less than 55, we use a simple length prefix
		return append(encodeLength(length, 0x80), []byte(str)...)
	} else {
		// For lengths >= 55, we need to handle longer length encoding
		lenOfLength := intsize(uint64(length))    // Compute size of length
		prefix := encodeLength(lenOfLength, 0xb7) // Adjust for offset 0xb7 for long strings
		lengthBytes := encodeLength(length, 0)    // Encode the actual length
		return append(append(prefix, lengthBytes...), []byte(str)...)
	}
}

func encodeByteSlice(data interface{}) []byte {
	t := reflect.ValueOf(data)
	byteSlice := t.Bytes()
	length := len(byteSlice)

	if length < 55 {
		// For lengths less than 55, we use a simple length prefix
		return append(encodeLength(length, 0x80), []byte(byteSlice)...)
	} else {
		// For lengths >= 55, we need to handle longer length encoding
		lenOfLength := intsize(uint64(length))
		prefix := encodeLength(lenOfLength, 0xb7)
		lengthBytes := encodeLength(length, 0)
		return append(append(prefix, lengthBytes...), []byte(byteSlice)...)
	}
}

func isByte(t reflect.Type) bool {
	return t.Kind() == reflect.Uint8
}
