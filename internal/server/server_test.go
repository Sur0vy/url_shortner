package server

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func Test_generateAddress(t *testing.T) {

	tests := []struct {
		name  string
		param string
		want  string
	}{
		{
			name:  "ok test 1",
			param: "1234",
			want:  ":1234",
		},
		{
			name:  "ok test 1",
			param: "",
			want:  ":",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := ":" + tt.param
			assert.Equal(t, tt.want, port) // меняем на функцию Equal из пакета assert
		})
	}
}
