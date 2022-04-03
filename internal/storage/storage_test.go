package storage

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMapStorage_Get(t *testing.T) {
	type fields struct {
		counter int
		data    map[int]URL
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MapStorage{
				counter: tt.fields.counter,
				data:    tt.fields.data,
			}
			got, err := s.Get(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapStorage_Insert(t *testing.T) {
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
			ms := NewMapStorage()
			sh := ms.Insert(tt.fields.data[tt.fields.counter].Full)
			//пока обработчик ошибок не предусмотрен, но над тестом стоит подумать
			if !tt.wantErr {
				assert.Equal(t, sh, tt.fields.data[tt.fields.counter].Short)
				assert.Equal(t, ms.counter, tt.fields.counter)
			}
		})
	}
}

func TestMapStorage_getFullURL(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MapStorage{
				counter: tt.fields.counter,
				data:    tt.fields.data,
			}
			got, err := s.getFullURL(tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFullURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFullURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMapStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MapStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMapStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMapStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
