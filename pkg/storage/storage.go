package storage

import (
	"fmt"
	"network/pkg/student"
)

type Storage struct {
	db map[string]*student.Student
}

func BuildStorage() Storage {
	return Storage{make(map[string]*student.Student)}
}

func (strg *Storage) Add(st *student.Student) {
	strg.db[st.GetName()] = st
}

func (strg *Storage) PrintAll() {
	fmt.Println("students from storage:")
	for _, s := range strg.db {
		fmt.Println(s.GetName(), s.GetAge(), s.GetGrade())
	}

}
