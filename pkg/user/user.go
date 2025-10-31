package user

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type User struct {
	Name    string `json:"Name"`
	Age     int    `json:"Age"`
	Friends []int  `json:"Friends"`
}

type UserDB struct {
	dumpFileName string
	dumpFilePtr  *os.File
	// uses to check is file has been modified previusly
	dumpFileHash [16]byte

	Last int          `json:"last"`
	Db   map[int]User `json:"Db"`
}

// Initialize new user db with empty values or from
// json file if it exist
func (udb *UserDB) Init(fname string) {
	fStat, err := os.Stat(fname)

	if os.IsNotExist(err) {
		// db json file not exist - set all to zero
		// and create empty file
		fmt.Printf("db file created at - %v\n", fname)
		udb.Db = make(map[int]User)
		udb.dumpFilePtr, err = os.Create(fname)
		if err != nil {
			log.Fatalf("error! can't create dump file %v, with - %v\n", fname, err)
		}

		udb.dumpFileHash = md5.Sum([]byte{})
	} else {
		// db json file exist - read and initialize
		// values from it
		fmt.Printf("UserDB file has found at - %v\n", fname)
		udb.dumpFilePtr, _ = os.OpenFile(fname, os.O_RDWR, 0)

		var fileContent []byte

		if fStat.Size() > 0 {
			fileContent, _ = io.ReadAll(udb.dumpFilePtr)
			err := json.Unmarshal(fileContent, udb)
			if err != nil {
				log.Fatalf("Init(): error parsing json - %v\n", err)
			}
		} else {
			udb.Db = make(map[int]User)
		}
		udb.dumpFileHash = md5.Sum(fileContent)
	}
}

func (udb UserDB) Close() {
	err := udb.dumpFilePtr.Close()
	if err != nil {
		log.Fatalf("error closing UserDB - %v\n", err)
	}
}

// Dump user db map[] and user id counter to disc
func (udb *UserDB) Dump() {
	out, _ := json.Marshal(udb)

	udb.dumpFileHash = md5.Sum(out)
	fmt.Printf("Dump(): new hash - %x\n", udb.dumpFileHash)

	// reduce file size to zero and set I/O offset at file begin
	udb.dumpFilePtr.Truncate(0)
	udb.dumpFilePtr.Seek(0, 0)

	_, err := udb.dumpFilePtr.WriteString(string(out))
	if err != nil {
		log.Fatalf("Dump(): %v\n", err)
	}
}

// Fetch check if file was modified (by calculating md5 sum of content)
// to prevent unnecessary json unmarshalling and map[int]user.User realloc,
// but it still need to read whole file from disc.
func (udb *UserDB) Fetch() {
	// return I/O offset to beginning of the file
	// after every r/w operation
	udb.dumpFilePtr.Seek(0, 0)

	fileContent, _ := io.ReadAll(udb.dumpFilePtr)

	tmp := md5.Sum(fileContent)
	fmt.Printf("Fetch():\tcur hash - %x\n\t\ttst hash - %x\n", udb.dumpFileHash, tmp)
	// hash of file content is similar, nothing to unmarshall
	if udb.dumpFileHash == tmp {
		fmt.Printf("Fetch(): nothing to fetch\n")
		udb.dumpFilePtr.Seek(0, 0)
		return
	}

	err := json.Unmarshal(fileContent, udb)
	if err != nil {
		log.Fatalf("Fetch(): error parsing json - %v\n", err)
	}

	udb.dumpFilePtr.Seek(0, 0)

	fmt.Printf("Fetch(): fetched\n")
}

func (udb *UserDB) AddUser(u User) int {
	udb.Fetch()

	udb.Db[udb.Last] = u
	tmp := udb.Last
	udb.Last++

	udb.Dump()

	return tmp
}

func (udb *UserDB) DeleteUser(id int) {
	udb.Fetch()

	delete(udb.Db, id)
	udb.Dump()
}

func (udb *UserDB) CheckUser(id int) error {
	udb.Fetch()

	_, found := udb.Db[id]
	if found {
		return nil
	} else {
		return errors.New("user not exist!")
	}
}

func (udb UserDB) GetUserName(id int) string {
	u, _ := udb.Db[id]
	return u.Name
}

func (udb UserDB) GetUserAge(id int) int {
	u, _ := udb.Db[id]
	return u.Age
}

func (udb UserDB) GetUserFriends(id int) []int {
	u, _ := udb.Db[id]
	return u.Friends
}

func (udb *UserDB) MakeFriends(source, target int) error {
	udb.Fetch()

	userSource, found := udb.Db[source]
	if !found {
		return errors.New("unable to make friends!")
	}

	userTarget, found := udb.Db[target]
	if !found {
		return errors.New("unable to make friends!")
	}

	userSource.Friends = append(userSource.Friends, target)
	userTarget.Friends = append(userTarget.Friends, source)

	delete(udb.Db, source)
	delete(udb.Db, target)

	udb.Db[source] = userSource
	udb.Db[target] = userTarget

	udb.Dump()

	return nil
}

func (udb *UserDB) UpdateUserAge(id, newAge int) error {
	udb.Fetch()

	if u, found := udb.Db[id]; found {
		delete(udb.Db, id)
		u.Age = newAge
		udb.Db[id] = u

		udb.Dump()

		return nil
	}

	return errors.New("User not exist!")
}
