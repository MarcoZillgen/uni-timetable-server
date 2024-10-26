package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	PORT := os.Getenv("API_PORT")
	r := mux.NewRouter()
	r.HandleFunc("/api/notion/data", getNotionData).Methods("GET")
	log.Println("Server is running on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, r))
}
