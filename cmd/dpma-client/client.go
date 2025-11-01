package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"network/internal/dpma"
)

func requestData(port int) (dpma.ResultT, error) {
	reqUrl := fmt.Sprintf("http://localhost:%d%v", port, "/")
	req, err := http.NewRequest(http.MethodGet, reqUrl, bytes.NewReader([]byte{}))
	if err != nil {
		return dpma.ResultT{}, err
	}
	defer req.Body.Close()

	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"Charset":      {"utf-8"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dpma.ResultT{}, err
	}

	fmt.Println(res.Status)
	resBody, _ := io.ReadAll(res.Body)

	var data dpma.ResultT
	if err := json.Unmarshal(resBody, &data); err != nil {
		return dpma.ResultT{}, err
	}

	return data, nil
}

func main() {
	var (
		action  string
		reqPort int
	)

	flag.IntVar(&reqPort, "port", 8585, "enter port")
	flag.Parse()

	fmt.Printf("request sending to port - %v\n", reqPort)

	for {
		fmt.Printf("type \"get\" to retrive data\n")
		fmt.Scan(&action)

		switch action {

		case "get":
			data, err := requestData(reqPort)
			if err != nil {
				fmt.Println(err)
			} else {
				for _, sms := range data.Data.SMS[0] {
					fmt.Printf("%v\n", sms)
				}
				fmt.Println()
				for _, mms := range data.Data.MMS[0] {
					fmt.Printf("%v\n", mms)
				}
			}

		default:
			fmt.Printf("unknown action - %v\n", action)
		}
	}
}
