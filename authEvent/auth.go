package authEvent

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
)

// AuthEvent gets the info about the UserCreation auth event
type AuthEvent struct {
	Email    string `json:"email"`
	Metadata struct {
		CreatedAt time.Time `json:"createdAt"`
	} `json:"metadata"`
	UID string `json:"uid"`
}

// Item represents an individual ShoppingListItem
type Item struct {
	ID string `json:"id,omitempty"`

	Name string `json:"name"`

	Amount float32 `json:"amount"`

	Unit string `json:"unit,omitempty"`
}

// List represents an array of ShoppingListItem
type List struct {
	Name string `json:"name"`

	Type string `json:"type"`

	Items []Item `json:"items,omitempty"`
}

// User struct for DB
type User struct {
	UID        string   `json:"uid"`
	Email      string   `json:"email"`
	FriendsUID []string `json:"friends"`
}

// UserCreated event is triggered when the user account is created
func UserCreated(ctx context.Context, e AuthEvent) error {
	log.Printf("UserCreated function triggered by the creation of user: %q", e.UID)
	log.Printf("Created at: %v", e.Metadata.CreatedAt)
	if e.Email != "" {
		log.Printf("Email: %q", e.Email)
	}
	conf := &firebase.Config{ProjectID: "collabshop19"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return fmt.Errorf("cannot get the firebase app: %v", err)
	}

	dbClient, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("cannot get the firestore client: %v", err)
	}
	defer dbClient.Close()

	user := User{
		UID:   e.UID,
		Email: e.Email,
	}

	list := List{
		Type: "shopping",
		Items: []Item{
			{
				ID:     "",
				Name:   "",
				Amount: 0.0,
				Unit:   "",
			},
		},
	}

	usersResult, err := dbClient.Collection("users").Doc(e.UID).Create(context.Background(), user)
	if err != nil {
		log.Fatalf("error writing to firestore: %v", err)
	}
	shoppingListRef, writeResult, err := dbClient.Collection("users").Doc(e.UID).Collection("lists").Add(ctx, list)
	if err != nil {
		log.Fatalf("error writing to firestore: %v", err)
	}

	log.Print(usersResult)
	log.Print(shoppingListRef, writeResult)
	return nil
}
