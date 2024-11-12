package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
)

func GetDefaultData(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	notionKey := os.Getenv("NOTION_KEY")
	dbId := os.Getenv("NOTION_DB_ID")

	dataHelper(w, notionKey, dbId)
}

func GetData(w http.ResponseWriter, r *http.Request) {
	notionKey := mux.Vars(r)["notionKey"]
	dbId := mux.Vars(r)["dbId"]
	dataHelper(w, notionKey, dbId)
}

type EventData struct {
	Title     string
	StartTime float64
	EndTime   float64
	Color     string
	EventType string
	Day       string // monday, tuesday, etc.
	Place     string
}

func dataHelper(w http.ResponseWriter, notionKey string, dbId string) {
	client := notionapi.NewClient(notionapi.Token(notionKey))
	db, err := client.Database.Query(context.Background(), notionapi.DatabaseID(dbId), nil)
	if err != nil {
		panic(err)
	}
	data := db.Results

	var eventDatas []EventData

	for _, page := range data {
		// check all attributes
		titleProp, ok := page.Properties["Name"].(*notionapi.TitleProperty)
		startProp, ok1 := page.Properties["Start"].(*notionapi.NumberProperty)
		endProp, ok2 := page.Properties["End"].(*notionapi.NumberProperty)
		dayProp, ok3 := page.Properties["Day"].(*notionapi.SelectProperty)
		colorProp, ok4 := page.Properties["Color"].(*notionapi.SelectProperty)
		typeProp, ok5 := page.Properties["Type"].(*notionapi.SelectProperty)
		placeProp, ok6 := page.Properties["Place"].(*notionapi.RichTextProperty)

		if !ok || !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
			fmt.Println("Missing property/properties in page: " + page.ID)
			fmt.Println(ok, ok1, ok2, ok3, ok4, ok5, ok6)
			fmt.Println("Skipping this page")
			fmt.Println()
			continue
		}

		event := EventData{
			Title:     titleProp.Title[0].PlainText,
			StartTime: startProp.Number,
			EndTime:   endProp.Number,
			Day:       dayProp.Select.Name,
			Color:     colorProp.Select.Name,
			EventType: typeProp.Select.Name,
			Place:     placeProp.RichText[0].PlainText,
		}

		eventDatas = append(eventDatas, event)
	}

	fmt.Println("Data fetched successfully")
	fmt.Println(eventDatas)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(eventDatas); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
