package service

import (
	"testing"
)

func Test_createCardNumber(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "generate",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createCardNumber(); !cardNumberIsValid(got) {
				t.Errorf("createCardNumber() generate invalid card: %v,", got)
			}
		})
	}
}
