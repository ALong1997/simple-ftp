package common

import (
	"unsafe"
)

const Address = "localhost:5900"

// string转[]byte
// 利用string本来的底层数组
func Str2sbyte(s string) (b []byte) {
	*(*string)(unsafe.Pointer(&b)) = s
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 2*unsafe.Sizeof(&b))) = len(s)
	return
}

// []byte转string
// 利用[]byte本来的底层数组
func Sbyte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
