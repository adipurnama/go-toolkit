package runtimekit_test

import (
	"testing"

	"github.com/adipurnama/go-toolkit/runtimekit"
)

func TestCallerName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"case 1",
			"runtimekit_test.TestCallerName.func1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runtimekit.CallerName(); got != tt.want {
				t.Errorf("CallerName() = %v, want %v", got, tt.want)
			}
		})
	}
}
