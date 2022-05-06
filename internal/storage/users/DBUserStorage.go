package users

type DBUserStorage struct {
	//current  string
	//fileName string
	//counter  int
	//Data     map[string]string
	//mtx      sync.RWMutex
}

func NewDBUserStorage() UserStorage {
	//dir, _ := os.Executable()
	//return &MapUserStorage{
	//	counter:  0,
	//	fileName: dir + UsersFileName,
	//	//fileName: "/Users/Sur0vy/Projects/url_shortner/" + UsersFileName,
	//	Data: make(map[string]string),
	//}
	return nil
}

func (u *DBUserStorage) Add() (string, string) {
	//u.mtx.Lock()
	//defer u.mtx.Unlock()
	//u.counter++
	//user := User + strconv.Itoa(u.counter)
	//hash := generateHash(user)
	//u.Data[hash] = user
	//u.writeToFile()
	//return user, hash
	return "", ""
}

func (u *DBUserStorage) GetUser(hash string) string {
	//u.mtx.RLock()
	//defer u.mtx.RUnlock()
	//user := u.Data[hash]
	//return user
	return ""
}

func (u *DBUserStorage) HasUser(hash string) bool {
	//u.mtx.RLock()
	//defer u.mtx.RUnlock()
	//if _, found := u.Data[hash]; found {
	//	return true
	//}
	//return false
	return false
}

func (u *DBUserStorage) LoadFromFile() error {
	//file, err := os.OpenFile(u.fileName, os.O_RDONLY|os.O_CREATE, 0777)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	//
	//scanner := bufio.NewScanner(file)
	//
	//if !scanner.Scan() {
	//	return scanner.Err()
	//}
	//data := scanner.Bytes()
	//
	//u.mtx.Lock()
	//defer u.mtx.Unlock()
	//err = json.Unmarshal(data, &u.Data)
	//if err != nil {
	//	return err
	//}
	//u.counter = len(u.Data)
	//return nil
	return nil
}

func (u *DBUserStorage) GetCount() int {
	//return u.counter
	return 0
}
