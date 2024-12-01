package subnet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSubnet(t *testing.T) {
	type args struct {
		clientIP string
		subnet   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty subnet - OK",
			args: args{
				clientIP: "",
				subnet:   "",
			},
			want: true,
		},
		{
			name: "Empty client IP - FAIL",
			args: args{
				clientIP: "",
				subnet:   "127.0.0.1/24",
			},
			want: false,
		},
		{
			name: "Corrupted subnet - FAIL",
			args: args{
				clientIP: "",
				subnet:   "127.0!_#.0.1/24",
			},
			want: false,
		},
		{
			name: "Not in subnet - FAIL",
			args: args{
				clientIP: "67.0.0.1",
				subnet:   "127.0.0.1/24",
			},
			want: false,
		},
		{
			name: "In subnet - OK",
			args: args{
				clientIP: "67.0.0.1",
				subnet:   "127.0.0.1/24",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateSubnet(tt.args.clientIP, tt.args.subnet)
			assert.Equal(t, tt.want, got)
		})
	}
}
