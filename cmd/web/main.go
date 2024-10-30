package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
)

type application struct {
	Session *scs.SessionManager
	DB      *sql.DB
	DSN     string
}

func main() {
	//set up an app config
	app := application{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5431 user=postgres password=postgres dbname=users sslmode=disable", "Postgres connection")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = conn
	//get a session manager
	app.Session = getSession()
	// get application routes
	mux := app.routes()
	// print out messages
	log.Println("Starting server on port 8080...")
	//start server
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
