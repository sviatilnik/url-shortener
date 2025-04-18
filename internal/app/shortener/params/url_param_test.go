package params

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLConfigParam_Validate(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "#1",
			value: "http://google.com",
			want:  true,
		},
		{
			name:  "#2",
			value: " ",
			want:  false,
		},
		{
			name:  "#3",
			value: "it's a valid url",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewURLConfigParam(tt.value)
			assert.Equal(t, tt.want, p.Validate())
		})
	}
}
