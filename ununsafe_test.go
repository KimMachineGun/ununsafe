package ununsafe_test

import (
	"math"
	"testing"

	"github.com/KimMachineGun/ununsafe"
	"github.com/stretchr/testify/assert"
)

func TestStringBytes(t *testing.T) {
	a := assert.New(t)

	tt := []struct {
		s string
		b []byte
	}{
		{"", nil},
		{"a", []byte{'a'}},
		{"ab", []byte{'a', 'b'}},
		{"abc", []byte{'a', 'b', 'c'}},
	}
	for _, tc := range tt {
		a.Equal(tc.b, ununsafe.StringToBytes(tc.s))
		a.Equal(tc.s, ununsafe.BytesToString(tc.b))
	}
}

func TestSameSize(t *testing.T) {
	a := assert.New(t)

	type A struct {
		a int64
		b uint32
		c [4]byte
	}
	type B [16]byte
	type C struct {
		a [4]byte
		b int32
		c uint64
	}

	tt := []struct {
		a        A
		b        B
		c        C
		expected []byte
	}{
		{
			a: A{
				a: 1,
				b: 2,
				c: [4]byte{3, 4, 5, 6},
			},
			b: B{
				1, 0, 0, 0, 0, 0, 0, 0,
				2, 0, 0, 0,
				3, 4, 5, 6,
			},
			c: C{
				a: [4]byte{1, 0, 0, 0},
				b: 0,
				c: 433757350042533890,
			},
			expected: []byte{
				0x1, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x2, 0x0, 0x0, 0x0,
				0x3, 0x4, 0x5, 0x6,
			},
		},
		{
			a: A{
				a: -1,
				b: math.MaxUint32,
				c: [4]byte{255, 255, 255, 255},
			},
			b: B{
				255, 255, 255, 255, 255, 255, 255, 255,
				255, 255, 255, 255,
				255, 255, 255, 255,
			},
			c: C{
				a: [4]byte{255, 255, 255, 255},
				b: -1,
				c: math.MaxUint64,
			},
			expected: []byte{
				0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff,
			},
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run("", func(t *testing.T) {
			a.Equal(tc.expected, ununsafe.ValueToSlice[A, byte](tc.a))
			a.Equal(tc.expected, ununsafe.ValueToSlice[B, byte](tc.b))
			a.Equal(tc.expected, ununsafe.ValueToSlice[C, byte](tc.c))
			a.Equal(tc.a, ununsafe.SliceToValue[byte, A](tc.expected))
			a.Equal(tc.b, ununsafe.SliceToValue[byte, B](tc.expected))
			a.Equal(tc.c, ununsafe.SliceToValue[byte, C](tc.expected))
		})
	}
}
