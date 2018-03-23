package app

// Endpoints.go contains the REST API endpoints for the service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//ReadPeople handles returning multiple people from the /people request
func (a *App) ReadPeople(w http.ResponseWriter, req *http.Request) {
	log.Printf("Got GET ALL")
	people, err := dbGetPeople(a.Database, 0, -1)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Could not get people: %v", err.Error())
		return
	}
	j, err := json.Marshal(people)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "could not marshal people: %v", err.Error())
	}
	fmt.Fprint(w, string(j))
}

// CreatePerson creates a new person in the database with ID n
func (a *App) CreatePerson(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	p := Person{}
	if vars["id"] != "" {
		log.Printf("Got POST ID %v", vars["id"])
		i, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(409)
			fmt.Fprintf(w, "Invalid ID")
			log.Printf("invalid ID passed: %v", err.Error())
			return
		}
		p.id = i
	} else {
		log.Printf("Got POST with NO ID")
		i, err := dbGetNextID(a.Database)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error getting next ID.")
			log.Printf("error getting next id: %v", err.Error())
			return
		}
		p.id = i
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	err := json.Unmarshal(buf.Bytes(), &p)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error Creating Person. Invalid input Data.")
		log.Printf("error unmarshalling data: %v", err.Error())
		return
	}

	err = p.dbCreatePerson(a.Database)
	if err != nil {
		w.WriteHeader(409)
		fmt.Fprintf(w, "Error Creating Person. ID Already Exists.")
		log.Printf("Error creating person: %v", err.Error())
		return
	}
	fmt.Fprintf(w, "Created Person with ID %v.", p.id)
}

// ReadPerson creates a new person in the database with ID n
func (a *App) ReadPerson(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	p := Person{}
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(409)
		fmt.Fprintf(w, "Invalid ID")
		log.Printf("invalid ID passed: %v", err.Error())
		return
	}
	err = p.dbGetPerson(a.Database, id)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Person not found.")
		log.Printf("could not marshal person: %v", err.Error())
		return
	}
	j, err := json.Marshal(p)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Could not format person.")
		log.Printf("could not marshal person: %v", err.Error())
		return
	}
	fmt.Fprint(w, string(j))
}

// UpdatePerson updates a person in the database with ID
func (a *App) UpdatePerson(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	log.Printf("Got UPDATE (%v) ID %v", req.Method, vars["id"])
	p := Person{}
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(409)
		fmt.Fprintf(w, "Invalid ID")
		log.Printf("invalid ID passed: %v", err.Error())
		return
	}
	p.id = id
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	err = json.Unmarshal(buf.Bytes(), &p)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error Updating Person. Invalid input Data.")
		log.Printf("error unmarshalling data: %v", err.Error())
		return
	}

	err = p.dbUpdatePerson(a.Database)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Person not found.")
		log.Printf("error creating person: %v", err.Error())
		return
	}
	fmt.Fprintf(w, "Updated Person with ID %v.", p.id)
}

// UpdatePatchPerson updates a person in the database with ID and partial input
func (a *App) UpdatePatchPerson(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	log.Printf("Got UPDATE (%v) ID %v", req.Method, vars["id"])
	p := Person{}
	Prev := Person{}
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(409)
		fmt.Fprintf(w, "Invalid ID")
		log.Printf("invalid ID passed: %v", err.Error())
		return
	}
	p.id = id
	Prev.dbGetPerson(a.Database, id)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	err = json.Unmarshal(buf.Bytes(), &p)
	if p.FirstName != "" {
		Prev.FirstName = p.FirstName
	}
	if p.LastName != "" {
		Prev.LastName = p.LastName
	}
	if p.Phone != "" {
		Prev.Phone = p.Phone
	}
	if p.Email != "" {
		Prev.Email = p.Email
	}
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error Updating Person. Invalid input Data.")
		log.Printf("error unmarshalling data: %v", err.Error())
		return
	}
	err = Prev.dbUpdatePerson(a.Database)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Person not found.")
		log.Printf("error creating person: %v", err.Error())
		return
	}
	fmt.Fprintf(w, "Updated Person with ID %v.", p.id)
}

// DeletePerson creates a new person in the database with ID n
func (a *App) DeletePerson(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	log.Printf("Got DELETE ID %v", vars["id"])
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(409)
		fmt.Fprintf(w, "Invalid ID")
		log.Printf("invalid ID passed: %v", err.Error())
		return
	}
	p := Person{id: id}
	err = p.dbDeletePerson(a.Database)
	if err != nil {
		fmt.Fprintf(w, "error deleting person: %v", err.Error())
		return
	}
	fmt.Fprintf(w, "Deleted Person with ID %v ", p.id)
}

// ImportCSV imports a CSV formatted list of entries into the database
func (a *App) ImportCSV(w http.ResponseWriter, req *http.Request) {
	log.Printf("Got POST to Import")
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	if buf.Len() == 0 {
		w.WriteHeader(409)
		fmt.Fprintf(w, "No Data.")
		log.Printf("No data provided on import")
		return
	}
	fmt.Printf("%v", buf.String())
	cr := csv.NewReader(buf)
	i := 0
	for {
		line, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v", line[0])
		id, err := dbGetNextID(a.Database)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error getting next ID.")
			log.Printf("error getting next id: %v", err.Error())
			return
		}
		if line[0] == "FirstName" {
			continue
		}
		p := Person{
			id:        id,
			FirstName: line[0],
			LastName:  line[1],
			Email:     line[2],
			Phone:     line[3],
		}
		p.dbCreatePerson(a.Database)
		i++
	}
	fmt.Fprintf(w, "Created %v entries.", i)
}

// ExportCSV exports a CSV formatted list of entries into the database
func (a *App) ExportCSV(w http.ResponseWriter, req *http.Request) {
	log.Printf("Got Export")
	people, err := dbGetPeople(a.Database, 0, -1)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Could not get people: %v", err.Error())
		return
	}
	buf := new(bytes.Buffer)
	cw := csv.NewWriter(buf)
	headers := people[0].GetHeaders()
	if err := cw.Write(headers); err != nil {
		log.Printf("Failed to write headers: %v", err.Error())
	}
	for _, person := range people {
		values := person.ToSlice()
		if err := cw.Write(values); err != nil {
			log.Printf("Failed to write values: %v", err.Error())
		}
	}
	cw.Flush()
	fmt.Fprintf(w, "%+v", buf.String())
}
