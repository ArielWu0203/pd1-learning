// main : if this file is an app(application).
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	
	// gorilla
	"github.com/gorilla/mux"

	// mongodb-go-driver
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

var DB_URI = "mongodb://localhost:27017"

type Problem struct {
	Pid			int 	"bson:'pid'"
	Title		string 	"bson:'title'"
	Solution	bool 	"bson:'solution'"
	Acceptance	float64 "bson:'acceptance'"
}

func ConnectToDB(URI, db_name, collection_name string)  (*mongo.Client, *mongo.Collection){
	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect to DB!!!")
	fmt.Println("client %T", client)

	collection := client.Database(db_name).Collection(collection_name)
	return client, collection
}

func DisconnectToDB( client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err!= nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to DB closed...")
}
func ExampleInsertDocs(collection *mongo.Collection) {
	// Create several documents in DB
	Doc1 := Problem{1, "Zuma Game", false, 40.2}
	Doc2 := Problem{2, "Word Subsets", true, 46.2}
	Doc3 := Problem{3, "Word Search", false, 32.9}
	Docs := []interface{}{Doc1, Doc2, Doc3}

	insertResult, err := collection.InsertMany(context.TODO(), Docs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Docs: ", insertResult.InsertedIDs)
}

func main() {

	// Connect to ExmpleDB and return client & collection
	// client, collection := ConnectToDB(DB_URI, "ExampleDB", "Problems")

	// Add some Docs
	// ExampleInsertDocs(collection)

	// Close the Connection
	// DisconnectToDB(client)

	// Create router
	router := mux.NewRouter()

	router.HandleFunc("/", hello_handler).Methods("GET")	
	router.HandleFunc("/api/problems", AllProblems_handler).Methods("GET")
	router.HandleFunc("/api/problems/{Pid:[0-9]+}", FindProblemWithPid_handler).Methods("GET")

	http.ListenAndServe("localhost:3000", router)
}

func hello_handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World")
}

func AllProblems_handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "All Problem")

	client, collection := ConnectToDB(DB_URI, "ExampleDB", "Problems")
	findOptions := options.Find()
	findOptions.SetLimit(10)
	var results []*Problem

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		var elem Problem
		err := cur.Decode(&elem)

		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)

	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())
	for _, r := range results {
		fmt.Fprintf(w, "Pid : %v\n", r.Pid)
		fmt.Fprintf(w, "Title : %v\n", r.Title)
		fmt.Fprintf(w, "Solution : %v\n", r.Solution)
		fmt.Fprintf(w, "Acceptance : %v\n", r.Acceptance)
		fmt.Fprintf(w, "---------------------------\n")
		
	}

	DisconnectToDB(client)
	
}

func FindProblemWithPid_handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	number, _ := strconv.Atoi(vars["Pid"])

	// Connect to DB
	client, collection := ConnectToDB(DB_URI, "ExampleDB", "Problems")
	var problem Problem
	filter := bson.M {"pid" : number}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&problem)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Problem: %v\n", problem)
	fmt.Fprintf(w, "Pid : %v\n", problem.Pid)
	fmt.Fprintf(w, "Title : %v\n", problem.Title)
	fmt.Fprintf(w, "Solution : %v\n", problem.Solution)
	fmt.Fprintf(w, "Acceptance : %v\n", problem.Acceptance)
	DisconnectToDB(client)
}
