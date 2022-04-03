package app

import (
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

func TestHandler_CreateShortURL(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: tt.fields.storage,
			}
			h.CreateShortURL(tt.args.c)
		})
	}
}

func TestHandler_GetURL(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: tt.fields.storage,
			}
			h.GetURL(tt.args.c)
		})
	}
}

func TestHandler_ResponseError(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: tt.fields.storage,
			}
			h.ResponseError(tt.args.c)
		})
	}
}

func TestNewHandler(t *testing.T) {
	type args struct {
		storage storage.Storage
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
