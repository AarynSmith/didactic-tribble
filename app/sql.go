package app

//SQL.go Contains the SQL queries used in the app

import (
	_ "github.com/mattn/go-sqlite3" // SQLITE3 driver for database/sql
)

const sqlTableCreate = `
CREATE TABLE IF NOT EXISTS people
(
id SERIAL,
fname TEXT NOT NULL,
lname TEXT NOT NULL,
email TEXT NOT NULL,
phone TEXT NOT NULL,
CONSTRAINT products_pkey PRIMARY KEY (id)
)`

const sqlTableClear = `
DELETE FROM people
`

const sqlReadPeople = `
SELECT id, fname, lname, email, phone 
FROM people 
LIMIT ? 
OFFSET ?
`

const sqlGetNextID = `
SELECT IFNULL(MAX(id),0)+1 FROM people
`

const sqlCreatePerson = `
INSERT INTO people (id, fname, lname, email, phone) 
VALUES (?, ?, ?, ?, ?)
`

const sqlReadPerson = `
SELECT id, fname, lname, email, phone FROM people WHERE id = ?
`

const sqlUpdatePerson = `
UPDATE people 
SET fname = ?, lname = ?, email = ?, phone = ?  
WHERE id = ?
`

const sqlDeletePerson = `
DELETE from people WHERE id = ?
`
