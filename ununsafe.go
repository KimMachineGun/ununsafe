package ununsafe

import (
	"fmt"
	"unsafe"
)

// StringToBytes converts string to []byte without copy.
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts []byte to string without copy.
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// SizeOf is a wrapper of unsafe.Sizeof.
func SizeOf[T any]() uint64 {
	var t T
	return uint64(unsafe.Sizeof(t))
}

func ValueToValue[From, To any](from From) To {
	fromSize := SizeOf[From]()
	toSize := SizeOf[To]()
	if fromSize != toSize {
		panic(fmt.Sprintf("ununsafe.ScalarToScalar: size mismatch %d != %d", fromSize, toSize))
	}
	return *(*To)(unsafe.Pointer(&from))
}

func SliceToValue[From, To any](arr []From) To {
	l := uint64(len(arr))
	fromSize := SizeOf[From]()
	toSize := SizeOf[To]()
	if l*fromSize != toSize {
		panic(fmt.Sprintf("ununsafe.VectorToScalar: size mismatch %d*%d != %d", l, fromSize, toSize))
	}
	return *(*To)(unsafe.Pointer(unsafe.SliceData(arr)))
}

func ValueToSlice[From, To any](from From) []To {
	fromSize := SizeOf[From]()
	toSize := SizeOf[To]()
	if fromSize%toSize != 0 {
		panic(fmt.Sprintf("ununsafe.ScalarToVector: size mismatch %d%%%d != 0", fromSize, toSize))
	}
	return unsafe.Slice(
		(*To)(unsafe.Pointer(&from)),
		fromSize/toSize,
	)
}

func SliceToSlice[From, To any](arr []From) []To {
	l := uint64(len(arr))
	fromSize := SizeOf[From]()
	toSize := SizeOf[To]()
	if l*fromSize%toSize != 0 {
		panic(fmt.Sprintf("ununsafe.VectorToVector: size mismatch %d*%d%%%d != 0", l, fromSize, toSize))
	}
	return unsafe.Slice(
		(*To)(unsafe.Pointer(unsafe.SliceData(arr))),
		l*fromSize/toSize,
	)
}
