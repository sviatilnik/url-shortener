package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultConfig_Get(t *testing.T) {

	store := make(map[string]interface{})
	store["key"] = "value"
	store["int"] = 5

	tests := []struct {
		name  string
		store map[string]interface{}
		key   string
		want  interface{}
	}{
		{
			name:  "#1",
			store: store,
			key:   "key",
			want:  "value",
		},
		{
			name:  "#2",
			store: store,
			key:   "test",
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewDefaultConfig(tt.store)

			assert.Equal(t, tt.want, c.Get(tt.key))
		})
	}
}

func TestDefaultConfig_Set(t *testing.T) {

	tests := []struct {
		name    string
		key     string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "#1",
			key:     "key",
			value:   "value",
			wantErr: false,
		},
		{
			name:    "#2",
			key:     "key",
			value:   1,
			wantErr: false,
		},
		{
			name:    "#3",
			key:     "key",
			value:   nil,
			wantErr: false,
		},
		{
			name:    "#4",
			key:     " ",
			value:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewDefaultConfig(nil)
			err := c.Set(tt.key, tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
