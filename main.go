package main

import (
	"log"

	"github.com/unixblackhole/didactic-tribble/app"
)

func main() {
	a := app.App{}
	err := a.Initialize("AddressBook.sqlitedb")
	if err != nil {
		log.Fatalf("Error Initializing: %v", err.Error())
	}
	a.Run(":3001")
}
