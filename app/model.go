package app

// Model.go contains the SQL Database model and connection information for the service.

import (
	"database/sql"
	"fmt"
)

// Person is an address book entry for a person.
type Person struct {
	id        int
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
	Phone     string `json:"Phone"`
}

// connectDatabase Creates our database connection
// An error is returned if ther is an issue creating the database or the table.
func connectDatabase(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		fmt.Printf("could not open database: %v", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not open database: %v", err.Error())
	}
	err = createTable(db)
	if err != nil {
		fmt.Printf("Could not create table: %v", err.Error())
		return nil, err
	}
	return db, nil
}

// createTable creates an empty sql table to store the data in.
func createTable(db *sql.DB) error {
	if _, err := db.Exec(sqlTableCreate); err != nil {
		return err
	}
	return nil
}

// clearTable deletes any data in the database.
func clearTable(db *sql.DB) error {
	if _, err := db.Exec(sqlTableClear); err != nil {
		return err
	}
	return nil
}

// dbGetNextID gets the highest used ID number in the database + 1
func dbGetNextID(db *sql.DB) (int, error) {
	row := db.QueryRow(sqlGetNextID)
	id := 0
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error getting row: %v", err.Error())
	}
	return id, nil
}

// dbGetPeople returns a slice of Person(s) and error.
// db is a Database connection, Start is the begining offset,
// and count is the number of records to return.
// Count of -1 returns all records
func dbGetPeople(db *sql.DB, start, count int) ([]Person, error) {
	rows, err := db.Query(sqlReadPeople, count, start)
	if err != nil {
		return nil, fmt.Errorf("error getting people: %v", err.Error())
	}
	People := []Person{}
	for rows.Next() {
		p := Person{}
		err = rows.Scan(&p.id, &p.FirstName, &p.LastName, &p.Email, &p.Phone)
		if err != nil {
			return nil, fmt.Errorf("error getting row: %v", err.Error())
		}
		People = append(People, p)
	}
	if len(People) == 0 {
		return nil, fmt.Errorf("no people returned")
	}
	return People, nil
}

// GetHeaders returns a person's headers (fields names)
func (p *Person) GetHeaders() []string {
	return []string{
		"FirstName",
		"LastName",
		"Email",
		"Phone",
	}
}

// ToSlice returns a person as a slice of fields
func (p *Person) ToSlice() []string {
	return []string{
		p.FirstName,
		p.LastName,
		p.Email,
		p.Phone,
	}
}

// dbCreatePerson Inserts a new person into the database.
// An error will be returned if the ID is already in use.
func (p *Person) dbCreatePerson(db *sql.DB) error {
	if _, err := db.Exec(sqlCreatePerson,
		p.id, p.FirstName, p.LastName, p.Email, p.Phone); err != nil {
		return err
	}
	return nil
}

// dbGetPerson Gets a specific person from the database.
// An error will be returned if there are no people in the database.
func (p *Person) dbGetPerson(db *sql.DB, id int) error {
	row := db.QueryRow(sqlReadPerson, id)
	err := row.Scan(&p.id, &p.FirstName, &p.LastName, &p.Email, &p.Phone)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return fmt.Errorf("ID not found")
		}
		return fmt.Errorf("error getting row: %v", err.Error())
	}
	return nil
}

// dbUpdatePerson Updates a specified person in the database.
// An error will be returned if the person is not already in the database.
func (p *Person) dbUpdatePerson(db *sql.DB) error {
	prev := Person{}
	err := prev.dbGetPerson(db, p.id)
	if err != nil {
		return err
	}
	if _, err := db.Exec(sqlUpdatePerson,
		p.FirstName, p.LastName, p.Email, p.Phone, p.id); err != nil {
		return err
	}
	return nil
}

// dbDeletePerson Deletes a specified person from the database.
// An error will be returned if the person is not already in the database.
func (p *Person) dbDeletePerson(db *sql.DB) error {
	if _, err := db.Exec(sqlDeletePerson,
		p.id); err != nil {
		return err
	}
	return nil
}
