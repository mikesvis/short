package main

import (
	_ "net/http/pprof"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_valueOrNA(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Value is string",
			args: args{
				value: "is_dummy_string",
			},
			want: "is_dummy_string",
		},
		{
			name: "Value is N/A",
			args: args{
				value: "",
			},
			want: "N/A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valueOrNA(tt.args.value)
			assert.Equal(t, tt.want, result)
		})
	}
}

func Test_buildInfo(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "No build values provided",
			want: "Build version: N/A\nBuild date: N/A\nBuild commit: N/A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildInfo()
			assert.Equal(t, tt.want, result)
		})
	}
}
