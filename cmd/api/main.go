package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"personal-projects/webapp/pkg/repository"
	"personal-projects/webapp/pkg/repository/dbrepo"
)

const port = 8090

type application struct {
	DSN      string
	DB       repository.DatabaseRepo
	Domain   string
	JWSecret string
}

func main() {
	var app application
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain for application")
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5431 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Posgtres connection")
	flag.StringVar(&app.JWSecret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "signing secret")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	log.Printf("Starting api on port %d", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
