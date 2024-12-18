package main

import (
	"os"
	"personal-projects/webapp/pkg/repository/dbrepo"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"
	app.Session = getSession()
	app.DB = &dbrepo.TestDBRepo{}
	os.Exit(m.Run())
}
