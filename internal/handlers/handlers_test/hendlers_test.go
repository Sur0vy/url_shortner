package hendlers_test

import (
	"bytes"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/handlers"
	"github.com/Sur0vy/url_shortner.git/internal/server"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testURL1     string = "http://www.blabla.net/blablabla"
	testURL1JSON string = "{\"url\":\"http://www.blabla.net/blablabla\"}"
	testURL2     string = "http://www.blabla.net/11111"
	testURL2JSON string = "{\"url\":\"http://www.blabla.net/11111\"}"
	testURL3     string = "http://www.rrr.com/wer/ggfsd"
	testURL4     string = "some text"
	testURL5     string = "//%2F1/1"
	testURL5Resp string = "http://%2F1/1"

	testShortURL1     string = "1"
	testShortURL1JSON string = "{\"result\":\"1\"}"
	testShortURL2     string = "2"
	testShortURL2JSON string = "{\"result\":\"2\"}"
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
			name: "Test POST correct",
			args: args{
				body:    testURL1,
				trueVal: true,
			},
			want: want{
				body: testShortURL1,
				code: http.StatusCreated,
			},
		},
		{
			name: "Test POST no entry",
			args: args{
				body:    testURL2,
				trueVal: false,
			},
			want: want{
				body: testShortURL2,
				code: http.StatusCreated,
			},
		},
	}
	config.Params = *config.SetupConfig()
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
				body = bytes.NewBuffer([]byte(config.Params.BaseURL + "/" + tt.want.body))
				assert.Equal(t, fmt.Sprint(w.Body), fmt.Sprint(body))
			} else {
				body = bytes.NewBuffer([]byte(config.Params.BaseURL + tt.want.body))
				assert.NotEqual(t, fmt.Sprint(w.Body), fmt.Sprint(body))
			}
		})
	}
}

func TestHandler_GetFullURL(t *testing.T) {
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
			name: "Test GET correct",
			fields: fields{
				body: []string{testURL1},
			},
			args: args{
				id: testShortURL1,
			},
			want: want{
				URL:  testURL1,
				code: http.StatusTemporaryRedirect,
			},
		},
		{
			name: "Test GET correct from many",
			fields: fields{
				body: []string{testURL1, testURL2, testURL3},
			},
			args: args{
				id: testShortURL2,
			},
			want: want{
				URL:  testURL2,
				code: http.StatusTemporaryRedirect,
			},
		},
		{
			name: "Test GET entry not found",
			fields: fields{
				body: []string{testURL1},
			},
			args: args{
				id: testShortURL2,
			},
			want: want{
				URL:  "",
				code: http.StatusNotFound,
			},
		},
		{
			name: "Test GET no entry",
			fields: fields{
				body: []string{},
			},
			args: args{
				id: testShortURL1,
			},
			want: want{
				URL:  "",
				code: http.StatusNotFound,
			},
		},
		{
			name: "Test GET without http",
			fields: fields{
				body: []string{testURL5},
			},
			args: args{
				id: testShortURL1,
			},
			want: want{
				URL:  testURL5Resp,
				code: http.StatusTemporaryRedirect,
			},
		},
		{
			name: "Test GET several of the same",
			fields: fields{
				body: []string{testURL1, testURL1, testURL5},
			},
			args: args{
				id: testShortURL1,
			},
			want: want{
				URL:  testURL1,
				code: http.StatusTemporaryRedirect,
			},
		},
	}

	config.Params = *config.SetupConfig()
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

func TestNewBaseHandler(t *testing.T) {
	tests := []struct {
		name string
		want *handlers.BaseHandler
	}{
		{
			name: "Test creating handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := handlers.NewBaseHandler(storage.NewMapStorage())
			assert.NotNil(t, ms)
		})
	}
}

func TestBaseHandler_GetShortURL(t *testing.T) {
	type fields struct {
		body []string
	}
	type args struct {
		body    string
		trueVal bool
	}
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Test getShortURL correct",
			fields: fields{
				body: []string{testURL1},
			},
			args: args{
				body:    testURL1JSON,
				trueVal: true,
			},
			want: want{
				body: testShortURL1JSON,
				code: http.StatusOK,
			},
		},
		{
			name: "Test getShortURL no entry",
			fields: fields{
				body: []string{testURL1},
			},
			args: args{
				body:    testURL2JSON,
				trueVal: true,
			},
			want: want{
				body: testShortURL2JSON,
				code: http.StatusCreated,
			},
		},
		{
			name: "Test getShortURL: wrong JSON",
			fields: fields{
				body: []string{testURL1, testURL2, testURL3, testURL4},
			},
			args: args{
				body:    testURL4,
				trueVal: false,
			},
			want: want{
				body: "",
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Test getShortURL correct, many entries in storage",
			fields: fields{
				body: []string{testURL1, testURL2, testURL3},
			},
			args: args{
				body:    testURL2JSON,
				trueVal: true,
			},
			want: want{
				body: testShortURL2JSON,
				code: http.StatusOK,
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

			//body := bytes.NewBuffer([]byte(tt.args.body))
			req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer([]byte(tt.args.body)))
			s.ServeHTTP(r, req)

			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, r.Code)
			if tt.args.trueVal {
				assert.Equal(t, r.Body, bytes.NewBuffer([]byte(tt.want.body)))
			} else {
				assert.Empty(t, r.Body)
			}
		})
	}
}
