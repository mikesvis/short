package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandkey(t *testing.T) {
	type want struct {
		typeOf  string
		len     int
		isEmpty bool
	}
	tests := []struct {
		name string
		arg  uint
		want want
	}{
		{
			name: "Rand key is string of length",
			arg:  5,
			want: want{
				typeOf:  "",
				len:     5,
				isEmpty: false,
			},
		}, {
			name: "Rand key is empty sting of zero length",
			arg:  0,
			want: want{
				typeOf:  "",
				len:     0,
				isEmpty: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			randKey := GetRandkey(tt.arg)
			assert.IsType(t, "", randKey)
			assert.Len(t, randKey, tt.want.len)
			if !tt.want.isEmpty {
				assert.NotEmpty(t, GetRandkey(tt.arg))
				return
			}

			assert.Empty(t, GetRandkey(tt.arg))
		})
	}
}
