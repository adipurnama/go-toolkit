package runtimekit_test

import (
	"testing"

	"github.com/adipurnama/go-toolkit/runtimekit"
)

func TestFunctionName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"case 1",
			"runtimekit_test.TestFunctionName.func1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runtimekit.FunctionName(); got != tt.want {
				t.Errorf("FunctionName() = %v, want %v", got, tt.want)
			}
		})
	}
}
