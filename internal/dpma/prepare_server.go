package dpma

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func handleConnection(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte("OK"))
	if err != nil {
		log.Println(err)
	}
}

func StartServer(port int) {
	router := chi.NewRouter()

	router.Get("/", handleConnection)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	fmt.Printf("Server listening at port - %v\n", port)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error running http server: %s\n", err)
		}
	}
}
