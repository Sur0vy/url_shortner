package users

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"strconv"
	"sync"
)

type MapUserStorage struct {
	current  string
	fileName string
	counter  int
	Data     map[string]string
	mtx      sync.RWMutex
}

func NewMapUserStorage() UserStorage {
	dir, _ := os.Executable()
	return &MapUserStorage{
		counter:  0,
		fileName: dir + UsersFileName,
		//fileName: "/Users/Sur0vy/Projects/url_shortner/" + UsersFileName,
		Data: make(map[string]string),
	}
}

func (u *MapUserStorage) Add(_ context.Context) (string, string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	u.counter++
	user := User + strconv.Itoa(u.counter)
	hash := GenerateHash(user)
	u.Data[hash] = user
	_ = u.writeToFile()
	return user, hash
}

func (u *MapUserStorage) GetUser(_ context.Context, hash string) string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	user := u.Data[hash]
	return user
}

func (u *MapUserStorage) LoadFromFile() error {
	file, err := os.OpenFile(u.fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return scanner.Err()
	}
	data := scanner.Bytes()

	u.mtx.Lock()
	defer u.mtx.Unlock()
	err = json.Unmarshal(data, &u.Data)
	if err != nil {
		return err
	}
	u.counter = len(u.Data)
	return nil
}

func (u *MapUserStorage) writeToFile() error {
	file, err := os.OpenFile(u.fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	data, err := json.Marshal(u.Data)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	// записываем URL в буфер
	if _, err := writer.Write(data); err != nil {
		return err
	}
	// записываем буфер в файл
	return writer.Flush()
}

func (u *MapUserStorage) GetCount(_ context.Context) int {
	return u.counter
}
