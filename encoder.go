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

// Calculate the size of an interface{} dynamically
func sizeOfInterface(data interface{}) uintptr {
	value := reflect.ValueOf(data)
	var size uintptr

	switch value.Kind() {
	case reflect.Slice:
		// Calculate the size of each element in the slice
		for i := 0; i < value.Len(); i++ {
			size += sizeOfInterface(value.Index(i).Interface())
		}
		// Add the overhead of the slice itself (slice header)
		size += unsafe.Sizeof(data)

	case reflect.Struct:
		// Calculate the size of each field in the struct
		for i := 0; i < value.NumField(); i++ {
			size += sizeOfInterface(value.Field(i).Interface())
		}

	case reflect.String:
		// The size of a string includes the size of the string header plus the length of the string
		size = unsafe.Sizeof("") + uintptr(len(data.(string)))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		// Use unsafe.Sizeof for primitive types
		size = unsafe.Sizeof(data)

	default:
		// If it's an unsupported type, use its unsafe size
		size = unsafe.Sizeof(data)
	}

	return size
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

func encodeLengthForLargeStructs(length int, lenOfLen int, offset byte) []byte {
	return []byte{byte(length) + byte(lenOfLen) + offset}
}

func encodeCustomSlice(e *Enc, data interface{}) []byte {
	var buffer bytes.Buffer
	slice := reflect.ValueOf(data)

	// Encode each element of the slice
	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i).Interface()
		if slice.Index(i).Kind() != reflect.Slice {
			// If element is not a slice, RLP-encode it directly
			buffer.Write(e.EncodeRLP(elem))
		} else {
			// Recursively encode if the element is a slice
			buffer.Write(encodeCustomSlice(e, elem))
		}
	}

	listSize := buffer.Len()

	if listSize < 56 {
		// List length is less than 56 bytes, use single byte prefix (0xc0 to 0xf7)
		return append(encodeLength(listSize, 0xc0), buffer.Bytes()...)
	} else {
		// List length is 56 bytes or more, use a multi-byte prefix (0xf8 to 0xff)
		encodedLength := encodeLength(listSize, 0x80) // encode the length as a binary
		prefixLength := len(encodedLength)            // number of bytes to store the length
		prefix := []byte{0xf7 + byte(prefixLength)}   // RLP large list prefix
		return append(append(prefix, encodedLength...), buffer.Bytes()...)
	}
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
