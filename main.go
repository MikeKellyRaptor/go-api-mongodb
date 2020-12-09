package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/mikekellyraptor/go-api-mongodb/utilities"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

var client = mongointerface.MongoConnect()

func main() {
	// MUX request routers
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/event/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/event/{id}", updateEvent).Methods("PATCH")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "You haven't made a correct request.")
	}

	collection := client.Database("event_handler_db").Collection("events")

	json.Unmarshal(reqBody, &newEvent)
	insertResult, err := collection.InsertOne(context.TODO(), newEvent)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Couldn't created data => %v", err)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(insertResult)
	}
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	filter := bson.D{{"id", eventID}}
	var result event

	collection := client.Database("event_handler_db").Collection("events")

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Couldn't find any data for ID %v => %v", eventID, err)
	} else {
		json.NewEncoder(w).Encode(result)
	}
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newEvent)

	filter := bson.D{{"id", eventID}}
	update := bson.D{
		{"$set", bson.D{
			{"Title", newEvent.Title},
			{"Description", newEvent.Description},
		}},
	}

	collection := client.Database("event_handler_db").Collection("events")

	updateresult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Couldn't update data for ID %v => %v", eventID, err)
		json.NewEncoder(w).Encode(newEvent)
	} else {
		json.NewEncoder(w).Encode(updateresult)
	}
}
