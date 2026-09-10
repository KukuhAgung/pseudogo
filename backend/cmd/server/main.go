package main

import (
	"log"
	"net/http"

	"pseudogo/internal/httpapi"
)

func main() {
	log.Println("Server jalan di 0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", httpapi.NewRouter()); err != nil {
		log.Fatal(err)
	}
}