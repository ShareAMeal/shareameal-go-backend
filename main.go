package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/macaron.v1"
	"log"
	"time"
)

type Association struct {
	MongoID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ID           int                `bson:"id" json:"id"`
	ContactEmail string             `json:"contact_email" bson:"contact_email"` // il faut bien que la 1e lettre soit en Majuscule sinon le champ est PRIVÉ !!!
	Name         string             `json:"name"`
	Location     string             `json:"location"`
	Description  string             `json:"description"`
	AdminId      int                `json:"admin" bson:"admin_id"` //	Admin        primitive.ObjectID `json:"admin" bson:"admin"`
	Phone        string             `json:"phone"`
}

type Event struct {
	ID            int       `json:"id" bson:"id"`
	Name          string    `json:"name" bson:"name"`
	StartDatetime time.Time `json:"start_datetime" bson:"start_datetime"`
	Active        bool      `json:"active" bson:"active"`
	Description   string    `json:"description" bson:"description"`
	OrganizerId   int       `json:"organizer" bson:"organizer_id"`
}

func main() {

	serveHttp(setupDb())

	/*err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")*/

}
func serveHttp(assoColl *mongo.Collection, eventColl *mongo.Collection) {
	m := macaron.Classic()

	m.Use(macaron.Renderer(macaron.RenderOptions{IndentJSON: true}))

	m.Use(func(ctx *macaron.Context) { // c'est un middleware
		log.Print("DNT: " + ctx.Req.Header.Get("dnt"))
		log.Print("User-Agent: " + ctx.Req.Header.Get("User-Agent"))
	})
	m.Use(func(ctx *macaron.Context) { // ça pourrait servir à faire de l'authentification
		username, password, ok := ctx.Req.Request.BasicAuth()
		if ok {
			log.Printf("%s attempted login with %s", username, password)
		} else {
			log.Print("not ok")
		}
	})

	m.Get("/asso/", func(ctx *macaron.Context) {
		ctx.JSON(200, listAssos(assoColl))
	})

	m.Get("/event/", func(ctx *macaron.Context) {
		ctx.JSON(200, listEvents(eventColl))
	})

	m.Get("/echo/:echostr", func(ctx *macaron.Context) {
		_, _ = ctx.Write([]byte(fmt.Sprintf("%v", ctx.AllParams())))
		ctx.Write([]byte("\nvous avez dit: " + ctx.Params(":echostr")))
	})

	m.Get("/echo/", func(ctx *macaron.Context) string { //https://go-macaron.com/middlewares/routing#named-parameters
		log.Printf("query: %v", ctx.QueryStrings("q"))
		//ctx.HTML(200,"vous avez dit: "+ctx.Params(":echostr"))
		return "return vous avez dit: " + ctx.Params(":echostr")
	})

	m.Run()
}

func setupDb() (*mongo.Collection, *mongo.Collection) {
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
	assoCollection := database.Collection("api_assocation")
	eventCollection := database.Collection("api_event")
	return assoCollection, eventCollection
}

func listAssos(assoCollection *mongo.Collection) []Association {
	defer timeTrack(time.Now(), "listAssos")
	var assos []Association
	findOptions := options.Find()
	findOptions.SetLimit(100)
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
	err = cur.Close(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v+", assos)
	return assos
}

func listEvents(eventCollection *mongo.Collection) []Event {
	var events []Event
	findOptions := options.Find()
	//findOptions.SetLimit(10)
	cur, err := eventCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Event
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		events = append(events, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v+", events)
	return events
}

func timeTrack(start time.Time, name string) { //https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
