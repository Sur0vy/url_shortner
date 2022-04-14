package storage

import (
	"bufio"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMapStorage_GetFullURL(t *testing.T) {
	type fields struct {
		counter int
		data    map[int]URL
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test get #1",
			fields: fields{
				counter: 1,
				data: map[int]URL{
					1: {
						Full:  "www.blabla.ru",
						Short: config.HTTP + config.HostAddr + ":" + config.HostPort + "/" + "1",
					},
				},
			},
			args: args{
				shortURL: "1",
			},
			want:    "www.blabla.ru",
			wantErr: false,
		},
		{
			name: "Test get #2",
			fields: fields{
				counter: 1,
				data: map[int]URL{
					1: {
						Full:  "www.blabla.ru",
						Short: "1",
					},
				},
			},
			args: args{
				shortURL: "2",
			},
			want:    "",
			wantErr: true,
		},
	}
	//config.HostAddr =
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MapStorage{
				Counter: tt.fields.counter,
				Data:    tt.fields.data,
			}
			fullURL, err := s.GetFullURL(tt.args.shortURL)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, fullURL, tt.want)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestMapStorage_InsertURL(t *testing.T) {
	type fields struct {
		counter int
		data    map[int]URL
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test insertr #1",
			fields: fields{
				counter: 1,
				data: map[int]URL{
					1: {
						Full:  "http://www.blabla.net/blablabla",
						Short: "http://localhost:8080/1",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMapStorage()
			sh := ms.InsertURL(tt.fields.data[tt.fields.counter].Full)
			//пока обработчик ошибок не предусмотрен, но над тестом стоит подумать
			if !tt.wantErr {
				assert.Equal(t, sh, tt.fields.data[tt.fields.counter].Short)
				assert.Equal(t, ms.Counter, tt.fields.counter)
			}
		})
	}
}

func TestNewMapStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MapStorage
	}{
		{
			name: "Test creating map storage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMapStorage()
			assert.NotNil(t, ms)
		})
	}
}

func TestMapStorage_GetShortURL(t *testing.T) {
	type fields struct {
		counter int
		data    map[int]URL
	}
	type args struct {
		fullURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test get #1",
			fields: fields{
				counter: 1,
				data: map[int]URL{
					1: {
						Full:  "www.blabla.ru",
						Short: "1",
					},
				},
			},
			args: args{
				fullURL: "www.blabla.ru",
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "Test get #2",
			fields: fields{
				counter: 1,
				data: map[int]URL{
					1: {
						Full:  "www.blabla.ru",
						Short: "1",
					},
				},
			},
			args: args{
				fullURL: "",
			},
			want:    "2",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MapStorage{
				Counter: tt.fields.counter,
				Data:    tt.fields.data,
			}
			fullURL, err := s.GetShortURL(tt.args.fullURL)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, fullURL.Short, tt.want)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestMapStorage_Load(t *testing.T) {
	type args struct {
		fileName string
		data     map[int]string
	}
	type fields struct {
		url map[int]URL
	}
	tests := []struct {
		name string
		args args
		want fields
	}{
		{
			name: "Test load from file #1",
			args: args{
				fileName: "\test.txt",
				data: map[int]string{
					1: `{"url":"http://www.werewrewr.com/f7","result":"http://localhost:8080/1"}`,
					2: `{"url": "http://www.werewrewr.com/f7/saf", "result": "http://localhost:8080/2"}`,
				},
			},
			want: fields{
				url: map[int]URL{
					1: {
						Full:  "http://www.werewrewr.com/f7",
						Short: "http://localhost:8080/1",
					},
					2: {
						Full:  "http://www.werewrewr.com/f7/saf",
						Short: "http://localhost:8080/2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := os.OpenFile(tt.args.fileName, os.O_RDWR|os.O_CREATE, 0777)
			defer os.Remove(tt.args.fileName)
			writer := bufio.NewWriter(file)
			for _, data := range tt.args.data {
				writer.WriteString(data + "\n")
			}
			writer.Flush()
			file.Close()
			ms := NewMapStorage()
			ms.Load(tt.args.fileName)

			for _, item := range tt.want.url {
				ShortURL, err := ms.GetShortURL(item.Full)
				assert.Nil(t, err)
				if err == nil {
					assert.Equal(t, item.Short, ShortURL.Short)
				}
			}
		})
	}
}

func TestMapStorage_addToFile(t *testing.T) {
	type args struct {
		fileName string
		data     map[int]URL
	}
	type fields struct {
		url map[int]URL
	}
	tests := []struct {
		name    string
		args    args
		want    fields
		wantErr bool
	}{
		{
			name: "Test add to file 2 item write/read",
			args: args{
				fileName: "\test.txt",
				data: map[int]URL{
					1: {
						Full:  "http://www.werewrewr.com/f7",
						Short: "http://localhost:8080/1",
					},
					2: {
						Full:  "http://www.werewrewr.com/f7/saf",
						Short: "http://localhost:8080/2",
					},
				},
			},
			want: fields{
				url: map[int]URL{
					1: {
						Full:  "http://www.werewrewr.com/f7",
						Short: "http://localhost:8080/1",
					},
					2: {
						Full:  "http://www.werewrewr.com/f7/saf",
						Short: "http://localhost:8080/2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test add to file no URL",
			args: args{
				fileName: "\test.txt",
				data: map[int]URL{
					1: {
						Full:  "http://www.werewrewr.com/f7/saf",
						Short: "http://localhost:8080/1",
					},
				},
			},
			want: fields{
				url: map[int]URL{
					1: {
						Full:  "http://www.werewrewr.com/f7",
						Short: "http://localhost:8080/1",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMapStorage()
			ms.Load(tt.args.fileName)
			defer os.Remove(tt.args.fileName)
			for _, item := range tt.args.data {
				ms.addToFile(&item)
			}

			ms2 := NewMapStorage()
			ms2.Load(tt.args.fileName)

			for _, data := range tt.want.url {
				ShortURL, err := ms2.GetShortURL(data.Full)
				if tt.wantErr == false {
					assert.Nil(t, err)
					if err == nil {
						assert.Equal(t, data.Short, ShortURL.Short)
					}
				} else {
					assert.NotNil(t, err)
				}
			}
			defer os.Remove(tt.args.fileName)
		})
	}
}
