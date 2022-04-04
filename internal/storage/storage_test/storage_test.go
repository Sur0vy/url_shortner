package storage_test

import (
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapStorage_Get(t *testing.T) {
	type fields struct {
		counter int
		data    map[int]storage.URL
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
				data: map[int]storage.URL{
					1: {
						Full:  "www.blabla.ru",
						Short: "1",
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
				data: map[int]storage.URL{
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage.MapStorage{
				Counter: tt.fields.counter,
				Data:    tt.fields.data,
			}
			fullURL, err := s.Get(tt.args.shortURL)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, fullURL, tt.want)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestMapStorage_Insert(t *testing.T) {
	type fields struct {
		counter int
		data    map[int]storage.URL
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
				data: map[int]storage.URL{
					1: {
						Full:  "www.blabla.ru",
						Short: "1",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := storage.NewMapStorage()
			sh := ms.Insert(tt.fields.data[tt.fields.counter].Full)
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
		want *storage.MapStorage
	}{
		{
			name: "Test creating map storage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := storage.NewMapStorage()
			assert.NotNil(t, ms)
		})
	}
}
