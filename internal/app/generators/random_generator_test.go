package generators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomGenerator_Get(t *testing.T) {
	tests := []struct {
		name    string
		len     uint
		str     string
		wantErr bool
	}{
		{
			name: "#1",
			len:  5,
			str:  "abc",
		},
		{
			name:    "#2",
			len:     10,
			str:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRandomGenerator(tt.len)
			got, err := r.Get(tt.str)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int(tt.len), len(got))
			}
		})
	}
}
