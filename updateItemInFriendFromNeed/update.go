package updateiteminfriendfromneed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

type Response struct {
	Body string `json:"body,omitempty"`
	Err  error  `json:"err,omitempty"`
}

type Request struct {
	RequesterUID    string `json:"uid"`
	RequesterListID string `json:"listid"`
	RequestedNeed   Need   `json:need"`
}

type Need struct {
	Name     string  `json:"name"`
	Quantity float32 `json:"quantity"`
	Type     string  `json:"type"`
	UOM      string  `json:"uom"`
}

type DataToBeWritten struct {
	OriginListUID string
	Type          string
}

// UpdateListToAddNeed Runs when the need quantity is updated
func UpdateListToAddNeed(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "collabshop19"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot get the firebase app: %v", err), http.StatusInternalServerError)
		return

	}

	req := &Request{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusInternalServerError)
		return
	}

	dbClient, err := app.Firestore(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating database client: %v", err), http.StatusInternalServerError)
		return
	}

	iter := dbClient.Collection("Users").Where("FriendsUID", "==", req.RequesterUID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		uid, _ := doc.DataAt("UID")

		needData := DataToBeWritten{
			OriginListUID: req.RequesterListID,
			Type:          "need",
		}
		listCreateResult, err := dbClient.Collection("Users").Doc(fmt.Sprint(uid)).Collection("Lists").Doc(req.RequesterListID).Create(ctx, needData)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating the list document: %v", err), http.StatusInternalServerError)
			return
		}
		addedItemRef, writeResult, err := dbClient.Collection("Users").Doc(fmt.Sprint(uid)).Collection("Lists").Doc(req.RequesterListID).Collection("Items").Add(ctx, req.RequestedNeed)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating the Item collection: %v", err), http.StatusInternalServerError)
			return
		}
		log.Print(listCreateResult)
		log.Print(addedItemRef, writeResult)
	}
}
