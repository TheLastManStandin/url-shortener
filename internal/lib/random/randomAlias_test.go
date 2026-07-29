package random

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNewRandomAlias(t *testing.T) {
	tests := []struct {
		name            string
		randAliasLength int
	}{
		{
			name:            "with 1",
			randAliasLength: 1,
		},
		{
			name:            "with 10",
			randAliasLength: 10,
		},
		{
			name:            "with 50",
			randAliasLength: 50,
		},
		{
			name:            "with 5",
			randAliasLength: 5,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomAlias(tt.randAliasLength)
			str2 := NewRandomAlias(tt.randAliasLength)

			assert.Equal(t, len(str1), tt.randAliasLength)
			assert.Equal(t, len(str2), tt.randAliasLength)

			assert.NotEqual(t, str1, str2)
		})
	}
}
