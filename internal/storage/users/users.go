package users

import (
	"context"
	"crypto/md5"
	"fmt"
)

const (
	User          string = "User_"
	UsersFileName string = "Users.txt"
)

type UserStorage interface {
	Add(ctx context.Context) (string, string)
	GetUser(ctx context.Context, hash string) string
	LoadFromFile() error
	GetCount(ctx context.Context) int
}

func GenerateHash(val string) string {
	h := md5.Sum([]byte(val))
	return fmt.Sprintf("%x", h)
}
