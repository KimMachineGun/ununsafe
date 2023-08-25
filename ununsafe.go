package ununsafe

import (
	"fmt"
	"unsafe"
)

func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func SizeOf[T any]() uint64 {
	var t T
	return uint64(unsafe.Sizeof(t))
}

func ScalarToScalar[From, To any](from From) To {
	fromSize := SizeOf[From]()
	toSize := SizeOf[To]()
	if fromSize != toSize {
		panic(fmt.Sprintf("ununsafe.ScalarToScalar: size mismatch %d != %d", fromSize, toSize))
	}
	return *(*To)(unsafe.Pointer(&from))
}

func VectorToScalar[From, To any](arr []From) To {
	l := uint64(len(arr))
	fromSize := SizeOf[From]()
	toSize := SizeOf[To]()
	if l*fromSize != toSize {
		panic(fmt.Sprintf("ununsafe.VectorToScalar: size mismatch %d*%d != %d", l, fromSize, toSize))
	}
	return *(*To)(unsafe.Pointer(unsafe.SliceData(arr)))
}

func ScalarToVector[From, To any](from From) []To {
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

func VectorToVector[From, To any](arr []From) []To {
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

func UpdateScalar[From, To any](from *From, fn func(*To)) {
	fn(ScalarToScalar[*From, *To](from))
}

func UpdateVector[From, To any](arr []From, fn func(To) To) []From {
	farr := VectorToVector[From, To](arr)
	for i := range farr {
		farr[i] = fn(farr[i])
	}
	return arr
}
