package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "#1 valid URL",
			url:  "http://www.google.com",
			want: true,
		},
		{
			name: "#2 invalid URL",
			url:  "www.google.com",
			want: false,
		},
		{
			name: "#3 invalid URL",
			url:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsURL(tt.url))
		})
	}
}
