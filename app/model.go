package app

import (
	"database/sql"
	"fmt"
)

// Person is an address book entry for a person.
type Person struct {
	ID        int    `json:"ID" csv:"ID"`
	FirstName string `json:"FirstName" csv:"FirstName"`
	LastName  string `json:"LastName" csv:"LastName"`
	Email     string `json:"Email" csv:"Email"`
	Phone     string `json:"Phone" csv:"Phone"`
}

const tableCreate = `CREATE TABLE IF NOT EXISTS people
(
id SERIAL,
fname TEXT NOT NULL,
lname TEXT NOT NULL,
email TEXT NOT NULL,
phone TEXT NOT NULL,
CONSTRAINT products_pkey PRIMARY KEY (id)
)`

func createTable(db *sql.DB) error {
	if _, err := db.Exec(tableCreate); err != nil {
		return err
	}
	return nil
}

func clearTable(db *sql.DB) error {
	if _, err := db.Exec("DELETE from people"); err != nil {
		return err
	}
	return nil
}

func getDBPeople(db *sql.DB, start, count int) ([]Person, error) {
	rows, err := db.Query("SELECT id, fname, lname, email, phone FROM people LIMIT ? OFFSET ?", count, start)
	if err != nil {
		return nil, fmt.Errorf("error getting people: %v", err.Error())
	}
	People := []Person{}
	for rows.Next() {
		p := Person{}
		err = rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email, &p.Phone)
		if err != nil {
			return nil, fmt.Errorf("error getting row: %v", err.Error())
		}
		People = append(People, p)
	}
	return People, nil
}

func getNextID(db *sql.DB) (int, error) {
	row := db.QueryRow("SELECT IFNULL(MAX(id),0)+1 FROM people")
	id := 0
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error getting row: %v", err.Error())
	}
	return id, nil
}

func (p *Person) dbGetPerson(db *sql.DB, id string) error {
	row := db.QueryRow("SELECT id, fname, lname, email, phone FROM people WHERE id = ?", id)
	err := row.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email, &p.Phone)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return fmt.Errorf("ID not found")
		}
		return fmt.Errorf("error getting row: %v", err.Error())
	}
	return nil
}

func (p *Person) dbUpdatePerson(db *sql.DB) error {
	if _, err := db.Exec("UPDATE people SET fname = ?, lname = ?, email = ?, phone = ?  WHERE id = ?",
		p.FirstName, p.LastName, p.Email, p.Phone, p.ID); err != nil {
		return err
	}
	return nil
}

func (p *Person) dbDeletePerson(db *sql.DB) error {
	if _, err := db.Exec("DELETE from people WHERE id = ?",
		p.ID); err != nil {
		return err
	}
	return nil
}

func (p *Person) dbCreatePerson(db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO people (id, fname, lname, email, phone) VALUES (?, ?, ?, ?, ?)",
		p.ID, p.FirstName, p.LastName, p.Email, p.Phone); err != nil {
		return err
	}
	return nil
}
