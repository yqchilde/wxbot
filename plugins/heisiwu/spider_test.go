package heisiwu

import "testing"

func Test_start(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"case1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start()
		})
	}
}
