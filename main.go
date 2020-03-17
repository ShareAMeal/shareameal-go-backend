package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"html"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Association struct {
	Contact_email string // il faut bien que la 1e lettre soit en Majuscule sinon le champ est PRIVÉ !!!
	Name          string
	Location      string
	Description   string
	//phone         string
	//admin         int
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://home.ribes.ovh:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	//on peut continuer

	database := client.Database("pweb")
	assoCollection := database.Collection("associations")
	/*assoTest := Association{
		//contact_email: "asso@insa-lyon.fr",
		Name:          "Asso 1",
		Location:      "20 avenue albert einsten",
		//description:   "dfgdgdfg",
		//phone:         "012345678987",
		//admin:         1,
	}*/
	/*insertResult,err := assoCollection.InsertOne(context.TODO(), bson.D{
		//{"contact_email","dfdg@dgfg.fg"},
		{"Name","Asso 2"},
		{"Location","21 avenue alber einstein"},
		//{"description","dfg fghf fg hfg f gh"},
		//{"phone","545654654"},
		//{"admin",1},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)*/
	assoTest := Association{
		Contact_email: "asso@insa-lyon.fr",
		Name:          "Asso 1",
		Location:      "20 avenue albert einsten",
		Description:   "dfgdgdfg",
		//phone:         "012345678987",
		//admin:         1,
	}
	insertResult, err := assoCollection.InsertOne(context.TODO(), assoTest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	var result Association
	err = assoCollection.FindOne(context.TODO(), bson.M{"name": "Asso 1"}).Decode(&result) //bson.D{{"Name", "Asso 1"}}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
	println("Name: " + result.Name)

	serveHttp(assoCollection)

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
func serveHttp(coll *mongo.Collection) {

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		assos := listAssos(coll)
		json, err := json.Marshal(assos)
		if err == nil {
			w.Write(json)
		} else {
			log.Fatal(err)
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
func listAssos(assoCollection *mongo.Collection) []Association {
	var assos []Association
	findOptions := options.Find()
	//findOptions.SetLimit(10)
	cur, err := assoCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem Association
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		assos = append(assos, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())
	return assos
}
