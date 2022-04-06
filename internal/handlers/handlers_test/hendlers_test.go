package hendlers_test

import (
	"bytes"
	"github.com/Sur0vy/url_shortner.git/internal/handlers"
	"github.com/Sur0vy/url_shortner.git/internal/server"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateShortURL(t *testing.T) {
	type args struct {
		body    string
		trueVal bool
	}
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test POST 1",
			args: args{
				body:    "http://www.blabla.net/blablabla",
				trueVal: true,
			},
			want: want{
				body: "http://localhost:8080/1",
				code: 201,
			},
		},
		{
			name: "Test POST 2",
			args: args{
				body:    "http://www.blabla.net/11111",
				trueVal: false,
			},
			want: want{
				body: "http://localhost:8080/2",
				code: 201,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.SetupServer()

			w := httptest.NewRecorder()

			body := bytes.NewBuffer([]byte(tt.args.body))
			req, err := http.NewRequest("POST", "/", body)

			s.ServeHTTP(w, req)

			assert.Nil(t, err)

			assert.Equal(t, w.Code, tt.want.code)

			if tt.args.trueVal {
				body = bytes.NewBuffer([]byte(tt.want.body))
				assert.Equal(t, w.Body, body)
			} else {
				body = bytes.NewBuffer([]byte(tt.want.body))
				assert.NotEqual(t, w.Body, body)
			}
		})
	}
}

func TestHandler_GetURL(t *testing.T) {
	type fields struct {
		body []string
	}
	type args struct {
		id string
	}
	type want struct {
		URL  string
		code int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Test GET 1",
			fields: fields{
				body: []string{"http://www.blabla.net/blablabla", "http://www.blabla.net/rrr", "ffff"},
			},
			args: args{
				id: "2",
			},
			want: want{
				URL:  "http://www.blabla.net/rrr",
				code: 307,
			},
		},
		{
			name: "Test GET 2",
			fields: fields{
				body: []string{"http://www.blabla.net/blablabla", "http://www.blabla.net/rrr", "ffff"},
			},
			args: args{
				id: "4",
			},
			want: want{
				URL:  "",
				code: 404,
			},
		},
		{
			name: "Test GET 3",
			fields: fields{
				body: []string{"http://www.blabla.net/blablabla", "http://www.blabla.net/rrr", "ffff"},
			},
			args: args{
				id: "1",
			},
			want: want{
				URL:  "http://www.blabla.net/blablabla",
				code: 307,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.SetupServer()
			w := httptest.NewRecorder()

			//заполним БД
			for _, v := range tt.fields.body {
				body := bytes.NewBuffer([]byte(v))
				req, _ := http.NewRequest("POST", "/", body)
				s.ServeHTTP(w, req)
			}

			r := httptest.NewRecorder()

			URL := "/" + tt.args.id
			req, err := http.NewRequest("GET", URL, nil)
			s.ServeHTTP(r, req)

			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, r.Code)

			URL = r.Header().Get("Location")
			assert.Equal(t, tt.want.URL, URL)
		})
	}
}

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name string
		want *handlers.Handler
	}{
		{
			name: "Test creating handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := handlers.NewHandler(storage.NewMapStorage())
			assert.NotNil(t, ms)
		})
	}
}
