package rfe

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFreeHostPort(t *testing.T) {
	t.Run("returned port is free", func(t *testing.T) {
		got := getFreeHostPort()

		l, err := net.Listen("tcp", fmt.Sprintf(":%d", got))
		defer l.Close()

		assert.NoError(t, err)
	})
}

func TestStringIsNotEmpty(t *testing.T) {
	tts := []struct {
		name string
		str  string
		want bool
	}{
		{name: "not empty string", str: "foo", want: true},
		{name: "empty string", str: "", want: false},
		{name: "empty string with whitespaces", str: " 	", want: false},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := stringIsNotEmpty(tt.str)

			assert.Equal(t, tt.want, got)
		})
	}

}
