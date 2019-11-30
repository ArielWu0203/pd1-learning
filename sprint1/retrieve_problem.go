package main

import (
	"context"
	"fmt"
	"log"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Problem struct {
	pid string
	title string 
	// solution int "bson:`Solution`"
	// Acceptence int "bson:`Acceptence`"
	// difficulty int "bson:`Diffculty`"
	// freqency int "bson:`Frequency`"
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect Success!!!")

	collection := client.Database("config").Collection("Problems")
	
	
}
