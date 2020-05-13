package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/macaron.v1"
	"log"
	"os"
	"time"
)

type Association struct {
	ID           int         `bson:"id" json:"id"`
	ContactEmail string      `json:"contact_email"`
	Name         string      `json:"name"`
	Location     string      `json:"location"`
	Description  string      `json:"description"`
	AdminId      int         `json:"admin"`
	Phone        interface{} `json:"phone"`
}

/* https://stackoverflow.com/a/33182597/12643213
func (asso *Association) Scan(res interface{}) error {
	ser, ok := res.(string)

	if !ok {
		return fmt.Errorf("wanted string, got %T instead", ser)
	}
	return json.Unmarshal([]byte(ser), asso)
}*/

type Event struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	StartDatetime time.Time `json:"start_datetime"`
	Active        bool      `json:"active"`
	Description   string    `json:"description"`
	OrganizerId   int       `json:"organizer"`
}

func main() {
	dbfile := os.Getenv("DB_FILE")
	if dbfile == "" {
		dbfile = "./db.sqlite3"
	}
	serveHttp(setupDb(dbfile))
}
func serveHttp(db *sql.DB) {
	m := macaron.Classic()
	m.Use(macaron.Renderer(macaron.RenderOptions{IndentJSON: true}))

	m.Get("/asso/", func(ctx *macaron.Context) {
		ctx.JSON(200, listAssos(db))
	})

	m.Get("/event/", func(ctx *macaron.Context) {
		ctx.JSON(200, listEvents(db))
	})

	/*m.Get("/echo/:echostr", func(ctx *macaron.Context) {
		_, _ = ctx.Write([]byte(fmt.Sprintf("%v", ctx.AllParams())))
		ctx.Write([]byte("\nvous avez dit: " + ctx.Params(":echostr")))
	})

	m.Get("/echo/", func(ctx *macaron.Context) string { //https://go-macaron.com/middlewares/routing#named-parameters
		log.Printf("query: %v", ctx.QueryStrings("q"))
		//ctx.HTML(200,"vous avez dit: "+ctx.Params(":echostr"))
		return "return vous avez dit: " + ctx.Params(":echostr")
	})*/

	m.Run()
}

func setupDb(databaseUrl string) *sql.DB {
	db, err := sql.Open("sqlite3", databaseUrl)
	handleError(err)
	err2 := db.PingContext(context.TODO())
	if err2 != nil {
		panic(err2)
	} else {
		fmt.Println("connexion SQL réussie")
	}
	return db
}

func listAssos(db *sql.DB) []Association {
	//rows, err := db.Query("select * from api_association;")
	rows, err := db.Query("select id, contact_email, name, location, description, phone, admin_id from api_association;")
	handleError(err)
	var assos []Association
	for rows.Next() {
		asso := Association{}
		err = rows.Scan(&asso.ID, &asso.ContactEmail, &asso.Name, &asso.Location, &asso.Description, &asso.Phone, &asso.AdminId)
		handleError(err)
		assos = append(assos, asso)
	}
	return assos
}

func listEvents(db *sql.DB) []Event {
	rows, err := db.Query("select id, name, start_datetime, active, description, organizer_id from api_event;")
	handleError(err)
	var events []Event
	for rows.Next() {
		event := Event{}
		handleError(rows.Scan(&event.ID, &event.Name, &event.StartDatetime, &event.Active, &event.Description, &event.OrganizerId))
		events = append(events, event)
	}
	return events
}

func timeTrack(start time.Time, name string) { //https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
