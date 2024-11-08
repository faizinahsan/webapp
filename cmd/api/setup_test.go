package main

import (
	"os"
	"personal-projects/webapp/pkg/repository/dbrepo"
	"testing"
)

var app application
var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE3MzA2OTEyNzUsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.buw8VEZLTcBG06gDgaNz-_3l3TNTsg8eETBM4PASKTs"

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "example.com"
	app.JWSecret = "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160"
	os.Exit(m.Run())
}
