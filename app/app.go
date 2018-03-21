package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// Import for SQLITE3
	_ "github.com/mattn/go-sqlite3"
)

const personInfoText = `ID:         %v
First Name: %v
Last Name:  %v
Phone Num:  %v
Email:      %v
`

// App contains an instanced router and assoicated backend database
type App struct {
	Router   *mux.Router
	Database *sql.DB
}

// Initialize creates our database instances
func (a *App) Initialize(dbname string) (err error) {
	a.Router = mux.NewRouter()
	a.Database, err = sql.Open("sqlite3", dbname)
	if err != nil {
		fmt.Printf("could not open database: %v", err.Error())
		return err
	}
	err = a.Database.Ping()
	if err != nil {
		return fmt.Errorf("could not open database: %v", err.Error())
	}
	err = createTable(a.Database)
	if err != nil {
		fmt.Printf("Could not create table: %v", err.Error())
		return err
	}
	return nil
}

// Index is the default response page for a GET to /
func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

// NotAllowed is the response for invalid requests
func NotAllowed(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(405)
	fmt.Fprintf(w, "Method Not Allowed.")
}

//GetPeople handles returning multiple people from the /people request
func (a *App) GetPeople(w http.ResponseWriter, req *http.Request) {
	people, err := getDBPeople(a.Database, 0, 10)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Could not get people: %v", err.Error())
		return
	}
	fmt.Fprintf(w, "Got %d people", len(people))
}

// GetPerson creates a new person in the database with ID n
func (a *App) GetPerson(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	p := Person{}
	err := p.dbGetPerson(a.Database, vars["id"])
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "error getting person: %v", err.Error())
		return
	}
	fmt.Fprintf(w, personInfoText, p.ID, p.FirstName, p.LastName, p.Phone, p.Email)
}

// // UpdateForm displays a form for updating a user in the database
// func (a *App) UpdateForm(w http.ResponseWriter, req *http.Request) {
// 	vars := mux.Vars(req)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		w.WriteHeader(404)
// 		fmt.Fprintf(w, "Cannot convert ID to int: %v", err.Error())
// 		return
// 	}
// 	p := Person{ID: id}
// 	err = p.dbGetPerson(a.Database, vars["id"])
// 	if err != nil {
// 		w.WriteHeader(404)
// 		fmt.Fprintf(w, "ID not found: %v", err.Error())
// 		return
// 	}
// 	fmt.Fprintf(w, `<html><body>
// <form action="/people/%v/update" method="post">
// <table>
// <tr><td>ID:</td><td><input type='text' value='%v' name='id' readonly /></td></tr>
// <tr><td>First Name:</td><td><input type='text' value='%v' name='fname' /></td></tr>
// <tr><td>Last Name:</td><td><input type='text' value='%v' name='lname' /></td></tr>
// <tr><td>Email:</td><td><input type='text' value='%v' name='email' /></td></tr>
// <tr><td>Phone:</td><td><input type='text' value='%v' name='phone' /></td></tr>
// <tr><td><input type="submit" value="Save"></td></tr>
// </table></body></html>
// 		`, p.ID, p.ID, p.FirstName, p.LastName, p.Email, p.Phone)
// }

// // UpdatePerson updates a person in the database with ID
// func (a *App) UpdatePerson(w http.ResponseWriter, req *http.Request) {
// 	vars := mux.Vars(req)
// 	fmt.Printf("%v", vars)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		fmt.Fprintf(w, "Cannot convert ID to int: \n%+v\n%v", vars["ID"], err.Error())
// 		return
// 	}
// 	p := Person{
// 		ID:        id,
// 		FirstName: vars["fname"],
// 		LastName:  vars["lname"],
// 		Email:     vars["email"],
// 		Phone:     vars["phone"],
// 	}
// 	err = p.dbUpdatePerson(a.Database)
// 	if err != nil {
// 		w.WriteHeader(409)
// 		fmt.Fprintf(w, "error creating person: %v", err.Error())
// 		return
// 	}

// 	fmt.Fprintf(w, "Entry Created:"+personInfoText, p.ID, p.FirstName, p.LastName, p.Phone, p.Email)
// }

// // UpdatePost updates a person in the database with ID
// func (a *App) UpdatePost(w http.ResponseWriter, req *http.Request) {
// 	vars := mux.Vars(req)
// 	fmt.Printf("%v", vars)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		fmt.Fprintf(w, "Cannot convert ID to int: \n%+v\n%v", vars["ID"], err.Error())
// 		return
// 	}
// 	p := Person{
// 		ID:        id,
// 		FirstName: req.PostForm["fname"][0],
// 		LastName:  req.PostForm["lname"][0],
// 		Email:     req.PostForm["email"][0],
// 		Phone:     req.PostForm["phone"][0],
// 	}
// 	err = p.dbDeletePerson(a.Database)
// 	if err != nil {
// 		w.WriteHeader(409)
// 		fmt.Fprintf(w, "error deleting person: %v", err.Error())
// 		return
// 	}
// 	err = p.dbCreatePerson(a.Database)
// 	if err != nil {
// 		w.WriteHeader(409)
// 		fmt.Fprintf(w, "error creating person: %v", err.Error())
// 		return
// 	}

// 	fmt.Fprintf(w, "Entry Created:"+personInfoText, p.ID, p.FirstName, p.LastName, p.Phone, p.Email)
// }

// CreatePerson creates a new person in the database with ID n
func (a *App) CreatePerson(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	id, err := strconv.Atoi(req.PostForm["id"][0])
	if err != nil {
		fmt.Fprintf(w, "Cannot convert ID to int: \n%+v\n%v", vars["ID"], err.Error())
		return
	}
	p := Person{
		ID:        id,
		FirstName: req.PostForm["fname"][0],
		LastName:  req.PostForm["lname"][0],
		Email:     req.PostForm["email"][0],
		Phone:     req.PostForm["phone"][0],
	}
	err = p.dbCreatePerson(a.Database)
	if err != nil {
		w.WriteHeader(409)
		fmt.Fprintf(w, "error creating person: %v", err.Error())
		return
	}

	fmt.Fprintf(w, "Entry Created:"+personInfoText, p.ID, p.FirstName, p.LastName, p.Phone, p.Email)
}

// CreateForm displays a new person form in the browser
func (a *App) CreateForm(w http.ResponseWriter, req *http.Request) {
	id, err := getNextID(a.Database)
	if err != nil {
		fmt.Fprintf(w, "Error getting next id: %v", err.Error())
		return
	}
	fmt.Fprintf(w, `<html><body>
<form action="/people/%v" method="post">
<table>
<tr><td>ID:</td><td><input type='text' value='%v' name='id' readonly /></td></tr>
<tr><td>First Name:</td><td><input type='text' name='fname' /></td></tr>
<tr><td>Last Name:</td><td><input type='text' name='lname' /></td></tr>
<tr><td>Email:</td><td><input type='text' name='email' /></td></tr>
<tr><td>Phone:</td><td><input type='text' name='phone' /></td></tr>
<tr><td><input type="submit" value="Save"></td></tr>
</table></body></html>
		`, id, id)
}

// DeletePerson creates a new person in the database with ID n
func (a *App) DeletePerson(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Cannot convert ID to int: \n%+v\n%v", vars["ID"], err.Error())
		return
	}
	p := Person{
		ID: id,
	}
	err = p.dbDeletePerson(a.Database)
	if err != nil {
		fmt.Fprintf(w, "error deleting person: %v", err.Error())
		return
	}

	fmt.Fprintf(w, "Entry ID %v Deleted", p.ID)
}

// Run starts an http listener on a specified address
func (a *App) Run(addr string) (err error) {
	a.addHandles()
	log.Printf("Listening on %v", addr)
	return http.ListenAndServe(addr, a.Router)
}

func (a *App) addHandles() {
	a.Router.HandleFunc("/", Index).Methods("GET")
	a.Router.HandleFunc("/", NotAllowed).Methods("POST", "PUT", "PATCH", "DELETE")
	a.Router.HandleFunc("/people", a.GetPeople).Methods("GET")
	a.Router.HandleFunc("/people", NotAllowed).Methods("POST", "PUT", "PATCH", "DELETE")
	a.Router.HandleFunc("/people/{id:[0-9]+}", a.GetPerson).Methods("GET")
	a.Router.HandleFunc("/people/{id:[0-9]+}", a.CreatePerson).Methods("POST")
	a.Router.HandleFunc("/people/{id:[0-9]+}", a.DeletePerson).Methods("DELETE")
	// a.Router.HandleFunc("/people/{id:[0-9]+}", a.UpdatePerson).Methods("PUT")
	a.Router.HandleFunc("/people/create", a.CreateForm).Methods("GET")
	// a.Router.HandleFunc("/people/{id:[0-9]+}/update", a.UpdateForm).Methods("GET")
	// a.Router.HandleFunc("/people/{id:[0-9]+}/update", a.UpdatePost).Methods("POST")
}
