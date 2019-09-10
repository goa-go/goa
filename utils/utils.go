package utils

import (
	"unsafe"
)

// Str2Bytes converts string to []byte.
func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Bytes2Str converts []byte to string.
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
