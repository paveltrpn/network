package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"network/pkg/citiesdb"
)

var (
	reqPort int = 3333
)

func makeAndSendReq(data []byte, method string, port int, url string) error {
	reqUrl := fmt.Sprintf("http://localhost:%d%v", port, url)
	req, err := http.NewRequest(method, reqUrl, bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer req.Body.Close()

	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"Charset":      {"utf-8"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(res.Status)
	resBody, _ := io.ReadAll(res.Body)
	fmt.Println(string(resBody))

	return nil
}

func requestList(port int) {
	err := makeAndSendReq([]byte{}, http.MethodGet, port, "/list")
	if err != nil {
		log.Fatalf("requestUpdateAge(): error sending request - %v\n", err)
	}
}

func requestGetById(id, port int) {
	err := makeAndSendReq([]byte{}, http.MethodGet, port, fmt.Sprintf("/get/%d", id))
	if err != nil {
		log.Fatalf("requestUpdateAge(): error sending request - %v\n", err)
	}
}

func requestAddCity(info citiesdb.CityInfo, port int) {
	reqJson, err := json.Marshal(info)
	if err != nil {
		log.Fatalf("requestAddCity(): error marshalling json - %v\n", err)
	}

	err = makeAndSendReq(reqJson, http.MethodPost, port, "/add")
	if err != nil {
		log.Fatalf("requestAddCity(): error sending request - %v\n", err)
	}
}

func requestDelete(id, port int) {
	err := makeAndSendReq([]byte{}, http.MethodDelete, port, fmt.Sprintf("/delete/%d", id))
	if err != nil {
		log.Fatalf("requestDelete(): error sending request - %v\n", err)
	}
}

func requestUpdatePopulation(id, count, port int) {
	params := make(map[string]int)
	params["id"] = id
	params["count"] = count

	reqJson, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("requestUpdatePopulation(): error marshalling json - %v\n", err)
	}

	err = makeAndSendReq(reqJson, http.MethodPost, port, "/update_pop")
	if err != nil {
		log.Fatalf("requestUpdatePopulation(): error sending request - %v\n", err)
	}

}

func requestByRegion(region string, port int) {
	err := makeAndSendReq([]byte{}, http.MethodGet, port, fmt.Sprintf("/get_by_region/%s", region))
	if err != nil {
		log.Fatalf("requestByRegion(): error sending request - %v\n", err)
	}
}

func requestByDistrict(district string, port int) {
	err := makeAndSendReq([]byte{}, http.MethodGet, port, fmt.Sprintf("/get_by_district/%s", district))
	if err != nil {
		log.Fatalf("requestByDistrict(): error sending request - %v\n", err)
	}
}

func main() {
	var (
		action string
	)

	flag.IntVar(&reqPort, "port", 3333, "enter port")
	flag.Parse()

	fmt.Printf("request sending to port - %v\n", reqPort)

	for {
		fmt.Printf("choose action - list get add delete update_pop by_region by_district\n")
		fmt.Scan(&action)

		switch action {
		case "list":
			requestList(reqPort)

		case "get":
			fmt.Printf("enter city id (number)\n")

			var id int
			fmt.Scanf("%v", &id)

			requestGetById(id, reqPort)

		case "add":
			fmt.Printf("enter name (string), region (string), district (string), population (int), foundation (int)\n")

			var (
				name       string
				region     string
				district   string
				population int
				foundation int
			)

			fmt.Scanf("%s %s %s %d %d", &name, &region, &district, &population, &foundation)

			requestAddCity(citiesdb.CityInfo{Id: 0, Name: name, Region: region, District: district, Population: population, Foundation: foundation}, reqPort)

		case "delete":
			fmt.Printf("ender id (number)\n")

			var id int

			fmt.Scanf("%v", &id)

			requestDelete(id, reqPort)

		case "update_pop":
			fmt.Printf("ender id (number) count (number)\n")

			var (
				id, count int
			)

			fmt.Scanf("%v %v", &id, &count)

			requestUpdatePopulation(id, count, reqPort)

		case "by_region":
			fmt.Printf("enter region (string)\n")

			var region string
			fmt.Scan(&region)

			requestByRegion(region, reqPort)

		case "by_district":
			fmt.Printf("enter district (string)\n")

			var district string
			fmt.Scan(&district)

			requestByDistrict(district, reqPort)

		default:
			fmt.Printf("unknown action - %v\n", action)
		}
	}
}
