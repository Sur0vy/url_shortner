package users

import (
	"crypto/md5"
	"fmt"
)

const (
	User          string = "User_"
	UsersFileName string = "Users.txt"
)

type UserStorage interface {
	Add() (string, string)
	GetUser(hash string) string
	HasUser(id string) bool
	LoadFromFile() error
	GetCount() int
}

func GenerateHash(val string) string {
	h := md5.Sum([]byte(val))
	return fmt.Sprintf("%x", h)
}
