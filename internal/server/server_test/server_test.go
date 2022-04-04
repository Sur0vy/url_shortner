package server_test

import (
	"github.com/Sur0vy/url_shortner.git/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupServer(t *testing.T) {
	tests := []struct {
		name string
		want *gin.Engine
	}{
		{
			name: "Test creating map storage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.SetupServer()
			assert.NotNil(t, s)
		})
	}
}
