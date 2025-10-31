package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"network/internal/dpma"
	"sync"

	"github.com/go-chi/chi"
)

const (
	host           = "http://localhost"
	srvPort        = 8585
	rscPort        = 8383
	assetsDir      = "../../assets/simulator/"
	servicesNumber = 7
)

func handleConnection(writer http.ResponseWriter, request *http.Request) {
	var (
		dataErr  error
		result   dpma.ResultT
		respCode int = http.StatusOK
		wg       sync.WaitGroup
	)

	wg.Add(servicesNumber)
	go func() {
		defer wg.Done()

		var err error

		result.Data.SMS, err = dpma.CollectSMSData(assetsDir + "sms.data")
		if err != nil {
			dataErr = err
		} else {
			log.Printf("sms data fetched!\n")
		}
	}()

	go func() {
		defer wg.Done()

		var err error

		result.Data.MMS, err = dpma.CollectMMSData(fmt.Sprintf("%v:%d%v", host, rscPort, "/mms"))
		if err != nil {
			dataErr = err
		} else {
			log.Printf("mms data fetched!\n")
		}
	}()

	go func() {
		defer wg.Done()

		var err error

		result.Data.VoiceCall, err = dpma.ReadVoiceCallData(assetsDir + "voice.data")
		if err != nil {
			dataErr = err
		} else {
			log.Printf("voice call data fetched!\n")
		}
	}()

	go func() {
		defer wg.Done()

		var err error

		result.Data.Email, err = dpma.CollectEmailData(assetsDir + "email.data")
		if err != nil {
			dataErr = err
		} else {
			log.Printf("email data fetched!\n")
		}
	}()

	go func() {
		defer wg.Done()

		var err error

		result.Data.Billing, err = dpma.ReadBillingData(assetsDir + "billing.data")
		if err != nil {
			dataErr = err
		} else {
			log.Printf("billing data fetched!\n")
		}
	}()

	go func() {
		defer wg.Done()

		var err error

		result.Data.Support, err = dpma.CollectSupportData(fmt.Sprintf("%v:%d%v", host, rscPort, "/support"))
		if err != nil {
			dataErr = err
		} else {
			log.Printf("support data fetched!\n")
		}
	}()

	go func() {
		defer wg.Done()

		var err error

		result.Data.Incidents, err = dpma.CollectIncidentData(fmt.Sprintf("%v:%d%v", host, rscPort, "/accendent"))
		if err != nil {
			dataErr = err
		} else {
			log.Printf("incidents data fetched!\n")
		}
	}()
	wg.Wait()

	if dataErr != nil {
		result.Error = dataErr.Error()
		respCode = http.StatusServiceUnavailable
	}

	respJson, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	writer.WriteHeader(respCode)

	_, err = writer.Write(respJson)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var (
		verbose     bool
		logFile     *os.File
		logFileName string = "server.log"
		logOutput   io.Writer
	)

	flag.BoolVar(&verbose, "v", false, "set verbose mode")
	flag.Parse()

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	if !verbose {
		logOutput = logFile
	} else {
		logOutput = io.MultiWriter(os.Stdout, logFile)
	}

	log.SetOutput(logOutput)

	router := chi.NewRouter()

	router.Get("/", handleConnection)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", srvPort),
		Handler: router,
	}

	log.Printf("Server listening at port - %v\n", srvPort)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error running http server: %s\n", err)
		}
	}
}
