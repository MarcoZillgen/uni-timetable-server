package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
)

func getNotionData(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	NOTION_KEY := os.Getenv("NOTION_KEY")
	DB_ID := os.Getenv("NOTION_DB_ID")

	client := notionapi.NewClient(notionapi.Token(NOTION_KEY))
	db, err := client.Database.Query(context.Background(), notionapi.DatabaseID(DB_ID), nil)
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

		if !ok || !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
			fmt.Println("Missing property/properties in page: " + page.ID)
			continue
		}

		event := EventData{
			Title:     titleProp.Title[0].PlainText,
			StartTime: startProp.Number,
			EndTime:   endProp.Number,
			Day:       dayProp.Select.Name,
			Color:     colorProp.Select.Name,
			EventType: typeProp.Select.Name,
		}

		eventDatas = append(eventDatas, event)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(eventDatas); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

type EventData struct {
	Title     string
	StartTime float64
	EndTime   float64
	Color     string
	EventType string
	Day       string // monday, tuesday, etc.
}
