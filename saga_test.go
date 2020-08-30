package saga

import "testing"

func TestStart(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Should Pass"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start()
		})
	}
}
