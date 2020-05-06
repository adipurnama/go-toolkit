package log

import (
	"testing"
)

func TestMaskPartial(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			"1 char",
			"a",
			"*",
		},
		{
			"2 char",
			"ab",
			"*b",
		},
		{
			"6 chars, starting 1 char & last 1 char remain unmasked",
			"abcdef",
			"a****f",
		},
		{
			"9 chars, starting 1 char & last 2 char remain unmasked",
			"123456789",
			"1******89",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskPartial(tt.args); got != tt.want {
				t.Errorf("MaskPartial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMask(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			"should mask string",
			"aaaaaaaa",
			"********",
		},
		{
			"no matter it's length",
			"aa",
			"**",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mask(tt.s); got != tt.want {
				t.Errorf("Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
