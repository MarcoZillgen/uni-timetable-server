package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MarcoZillgen/uni_plan/notion"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// load data
	godotenv.Load()
	PORT := os.Getenv("API_PORT")

	// api router
	r := mux.NewRouter()
	r.HandleFunc("/api/notion/data", notion.GetDefaultData).Methods("GET")
	r.HandleFunc("/api/notion/data", notion.GetData).Methods("GET").Queries("notionKey", "{notionKey}", "dbID", "{dbID}")
	log.Println("Server is running on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, r))
}
