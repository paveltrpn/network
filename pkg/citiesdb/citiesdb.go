package citiesdb

import (
	"bufio"
	"container/list"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CityInfo struct {
	Id         int    `json:"Id"`
	Name       string `json:"Name"`
	Region     string `json:"Region"`
	District   string `json:"District"`
	Population int    `json:"Population"`
	Foundation int    `json:"Foundation"`
}

type CitiesDB struct {
	dbFileName string
	dbFilePtr  *os.File

	Db *list.List
}

func (db *CitiesDB) Init(fname string) {
	var err error

	db.dbFileName = fname
	db.Db = list.New()

	fStat, err := os.Stat(fname)

	if os.IsNotExist(err) {
		// db json file not exist - set all to zero
		// and create empty file
		fmt.Printf("db file created at - %v\n", fname)

		db.dbFilePtr, err = os.Create(fname)
		if err != nil {
			log.Fatalf("error! can't create dump file %v, with - %v\n", fname, err)
		}
	} else {
		// db json file exist - read and initialize
		// values from it
		fmt.Printf("db file has found at - %v\n", fname)
		db.dbFilePtr, err = os.OpenFile(fname, os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err.Error())
		}

		if fStat.Size() > 0 {
			scanner := bufio.NewScanner(db.dbFilePtr)
			for scanner.Scan() {
				var (
					id                     int
					name, region, district string
					population, foundation int
				)

				line := strings.Split(scanner.Text(), ",")

				id, _ = strconv.Atoi(line[0])
				name = line[1]
				region = line[2]
				district = line[3]
				population, _ = strconv.Atoi(line[4])
				foundation, _ = strconv.Atoi(line[5])

				info := CityInfo{Id: id, Name: name, Region: region,
					District: district, Population: population, Foundation: foundation}

				db.Db.PushBack(&info)
			}
		}
	}
}

func (db *CitiesDB) Close() {
	db.dbFilePtr.Close()
}

// Dump all data to disc
func (db *CitiesDB) Dump() {
	// reduce file size to zero and set I/O offset at file begin
	db.dbFilePtr.Truncate(0)
	db.dbFilePtr.Seek(0, 0)

	w := bufio.NewWriter(db.dbFilePtr)

	for e := db.Db.Front(); e != nil; e = e.Next() {
		info := e.Value.(*CityInfo)
		_, err := w.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v\n",
			info.Id, info.Name, info.Region, info.District, info.Population, info.Foundation))
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if err := w.Flush(); err != nil {
		log.Fatal(err.Error())
	}
}

func (db *CitiesDB) GetById(id int) (CityInfo, error) {
	var rt CityInfo

	for e := db.Db.Front(); e != nil; e = e.Next() {
		if e.Value.(*CityInfo).Id == id {
			rt = *e.Value.(*CityInfo)
			return rt, nil
		}
	}
	return CityInfo{}, errors.New("not found")
}

func (db *CitiesDB) Add(infoStruct CityInfo) {
	// Make city id from two first bytes of md5 hash from string
	// formed by Name, District and Region fields of CityInfo struct.
	// WARN: may cause id collision!
	toHash := fmt.Sprintf("%v%v%v", infoStruct.Name, infoStruct.District, infoStruct.Region)
	md5Hash := md5.Sum([]byte(toHash))
	infoStruct.Id = int(binary.BigEndian.Uint16(md5Hash[0:2]))

	db.Db.PushBack(&infoStruct)
}

func (db *CitiesDB) Delete(id int) error {
	for e := db.Db.Front(); e != nil; e = e.Next() {
		if e.Value.(*CityInfo).Id == id {
			db.Db.Remove(e)
			return nil
		}
	}

	return errors.New("not found")
}

func (db *CitiesDB) UpdatePopulationInfo(id, newPopulation int) error {
	for e := db.Db.Front(); e != nil; e = e.Next() {
		if e.Value.(*CityInfo).Id == id {
			e.Value.(*CityInfo).Population = newPopulation
			return nil
		}
	}

	return errors.New("not found")
}

func (db *CitiesDB) GetByRegion(region string) (rt []CityInfo) {
	for e := db.Db.Front(); e != nil; e = e.Next() {
		if e.Value.(*CityInfo).Region == region {
			rt = append(rt, *e.Value.(*CityInfo))
		}
	}
	return rt
}

func (db *CitiesDB) GetByDistrict(district string) (rt []CityInfo) {
	for e := db.Db.Front(); e != nil; e = e.Next() {
		if e.Value.(*CityInfo).District == district {
			rt = append(rt, *e.Value.(*CityInfo))
		}
	}
	return rt
}

func (db *CitiesDB) GetByPopulation(from, to int) (rt []CityInfo) {
	for e := db.Db.Front(); e != nil; e = e.Next() {
		if (e.Value.(*CityInfo).Population > from) && (e.Value.(*CityInfo).Population < to) {
			rt = append(rt, *e.Value.(*CityInfo))
		}
	}
	return rt
}

func (db *CitiesDB) GetByFoundation(from, to int) (rt []CityInfo) {
	for e := db.Db.Front(); e != nil; e = e.Next() {
		if (e.Value.(*CityInfo).Foundation > from) && (e.Value.(*CityInfo).Foundation < to) {
			rt = append(rt, *e.Value.(*CityInfo))
		}
	}
	return rt
}
