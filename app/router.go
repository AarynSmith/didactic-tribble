package app

// Router.go contains combines the router and the database model to form the application.

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// App contains an instanced router and assoicated backend database
type App struct {
	Router   *mux.Router
	Database *sql.DB
}

// Initialize creates our database instances
func (a *App) Initialize(dbname string) (err error) {
	a.Router = mux.NewRouter()
	a.Database, err = connectDatabase(dbname)
	if err != nil {
		return fmt.Errorf("could not initialize: %v", err.Error())
	}
	return nil
}

// Run starts an http listener on a specified address
func (a *App) Run(addr string) (err error) {
	a.addHandles()
	log.Printf("Listening on %v", addr)
	return http.ListenAndServe(addr, a.Router)
}

// addHanles assings handler functions to the various methods and endpoints.
func (a *App) addHandles() {
	a.Router.HandleFunc("/people", a.ReadPeople).Methods("GET")
	a.Router.HandleFunc("/person", a.CreatePerson).Methods("POST")
	a.Router.HandleFunc("/person/{id:[0-9]+}", a.CreatePerson).Methods("POST")
	a.Router.HandleFunc("/person/{id:[0-9]+}", a.ReadPerson).Methods("GET")
	a.Router.HandleFunc("/person/{id:[0-9]+}", a.UpdatePerson).Methods("PUT")
	a.Router.HandleFunc("/person/{id:[0-9]+}", a.UpdatePatchPerson).Methods("PATCH")
	a.Router.HandleFunc("/person/{id:[0-9]+}", a.DeletePerson).Methods("DELETE")
	a.Router.HandleFunc("/import", a.ImportCSV).Methods("POST")
	a.Router.HandleFunc("/export", a.ExportCSV).Methods("GET")

}
