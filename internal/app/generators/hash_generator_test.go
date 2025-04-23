package generators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashGenerator_Get(t *testing.T) {

	tests := []struct {
		name    string
		len     uint
		str     string
		want    string
		wantErr bool
	}{
		{
			name: "#1",
			len:  10,
			str:  "hello world",
			want: "5eb63bbbe0",
		},
		{
			name:    "#2",
			len:     3,
			str:     " ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewHashGenerator(tt.len)
			got, err := g.Get(tt.str)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
