package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"network/pkg/citiesdb"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

var dbInstance citiesdb.CitiesDB

func handleListDb(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("list cities db\n")

	var rt []string

	for e := dbInstance.Db.Front(); e != nil; e = e.Next() {
		foo := e.Value.(*citiesdb.CityInfo)
		rt = append(rt, fmt.Sprintf("id - %v name - %v, region - %v, district - %v, population - %v, foundation - %v\n",
			foo.Id, foo.Name, foo.Region, foo.District, foo.Population, foo.Foundation))
	}

	_, err := writer.Write([]byte(strings.Join(rt, "")))
	if err != nil {
		log.Println(err)
	}
}

func handleGetById(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("get by id\n")

	cityId, _ := strconv.Atoi(chi.URLParam(request, "id"))

	cInfoStruct, err := dbInstance.GetById(cityId)

	var infoString string
	if err == nil {
		writer.WriteHeader(200)
		infoString = fmt.Sprintf("id - %v name - %v, region - %v, district - %v, population - %v, foundation - %v\n",
			cInfoStruct.Id, cInfoStruct.Name, cInfoStruct.Region, cInfoStruct.District, cInfoStruct.Population, cInfoStruct.Foundation)
	} else {
		writer.WriteHeader(503)
		infoString = err.Error()
	}

	_, err = writer.Write([]byte(infoString))
	if err != nil {
		log.Println(err)
	}
}

func handleAddCity(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("add city info\n")

	reqBody, _ := io.ReadAll(request.Body)

	var city citiesdb.CityInfo

	json.Unmarshal(reqBody, &city)

	dbInstance.Add(city)

	writer.WriteHeader(201)
	_, err := writer.Write([]byte("proceesed"))
	if err != nil {
		log.Println(err)
	}
}

func handleDelete(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("delete\n")

	cityId, _ := strconv.Atoi(chi.URLParam(request, "id"))

	err := dbInstance.Delete(cityId)

	var infoString string
	if err == nil {
		writer.WriteHeader(200)
		infoString = "delete success"
	} else {
		writer.WriteHeader(503)
		infoString = err.Error()
	}

	_, err = writer.Write([]byte(infoString))
	if err != nil {
		log.Println(err)
	}
}

func handleUpdatePopulation(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("update population\n")

	reqBody, _ := io.ReadAll(request.Body)

	var params map[string]int
	json.Unmarshal(reqBody, &params)

	var respString string
	if err := dbInstance.UpdatePopulationInfo(params["id"], params["count"]); err == nil {
		writer.WriteHeader(200)
		respString = "success"
	} else {
		writer.WriteHeader(503)
		respString = err.Error()
	}

	_, err := writer.Write([]byte(respString))
	if err != nil {
		log.Println(err)
	}
}

func handleGetByRegion(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("get by region\n")

	var rt []string

	info := dbInstance.GetByRegion(chi.URLParam(request, "region"))

	for _, foo := range info {
		rt = append(rt, fmt.Sprintf("id - %v name - %v, region - %v, district - %v, population - %v, foundation - %v\n",
			foo.Id, foo.Name, foo.Region, foo.District, foo.Population, foo.Foundation))
	}

	_, err := writer.Write([]byte(strings.Join(rt, "")))
	if err != nil {
		log.Println(err)
	}
}

func handleGetByDistrict(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("get by region\n")

	var rt []string

	info := dbInstance.GetByDistrict(chi.URLParam(request, "district"))

	for _, foo := range info {
		rt = append(rt, fmt.Sprintf("id - %v name - %v, region - %v, district - %v, population - %v, foundation - %v\n",
			foo.Id, foo.Name, foo.Region, foo.District, foo.Population, foo.Foundation))
	}

	_, err := writer.Write([]byte(strings.Join(rt, "")))
	if err != nil {
		log.Println(err)
	}
}

func handleGetByPopulation(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("get by population\n")

	reqBody, _ := io.ReadAll(request.Body)

	var params map[string]int
	json.Unmarshal(reqBody, &params)

	var rt []string
	info := dbInstance.GetByPopulation(params["from"], params["to"])
	for _, foo := range info {
		rt = append(rt, fmt.Sprintf("id - %v name - %v, region - %v, district - %v, population - %v, foundation - %v\n",
			foo.Id, foo.Name, foo.Region, foo.District, foo.Population, foo.Foundation))
	}

	writer.WriteHeader(200)
	_, err := writer.Write([]byte(strings.Join(rt, "")))
	if err != nil {
		log.Println(err)
	}
}

func handleGetByFoundation(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("get by population\n")

	reqBody, _ := io.ReadAll(request.Body)

	var params map[string]int
	json.Unmarshal(reqBody, &params)

	var rt []string
	info := dbInstance.GetByFoundation(params["from"], params["to"])
	for _, foo := range info {
		rt = append(rt, fmt.Sprintf("id - %v name - %v, region - %v, district - %v, population - %v, foundation - %v\n",
			foo.Id, foo.Name, foo.Region, foo.District, foo.Population, foo.Foundation))
	}

	writer.WriteHeader(200)
	_, err := writer.Write([]byte(strings.Join(rt, "")))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var (
		port       int
		dbFileName string
	)

	flag.StringVar(&dbFileName, "db", "db.csv", "enter file name to store data")
	flag.IntVar(&port, "port", 3333, "enter port")
	flag.Parse()

	dbInstance.Init(dbFileName)
	defer dbInstance.Close()

	router := chi.NewRouter()
	router.Get("/list", handleListDb)
	router.Get("/get/{id}", handleGetById)
	router.Post("/add", handleAddCity)
	router.Delete("/delete/{id}", handleDelete)
	router.Post("/update_pop", handleUpdatePopulation)
	router.Get("/get_by_region/{region}", handleGetByRegion)
	router.Get("/get_by_district/{district}", handleGetByDistrict)
	router.Get("/get_by_population", handleGetByPopulation)
	router.Get("/get_by_foundation", handleGetByFoundation)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("server started at port - %v\n", port)

	sig := <-done
	log.Printf("server stopped by signal - %v\n", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		dbInstance.Dump()
		fmt.Printf("data base dumped to - %v\n", dbFileName)
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	log.Print("server exited properly")
}
