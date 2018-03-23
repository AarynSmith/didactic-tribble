# didactic-tribble
REST API for an Address Book service

Default listening port is :3001

Methods available are:

* GET:
  * /people: Lists all people in the database.
  * /person/{id}: Gets a specific person by ID.
  * /export: Returns a CSV formated file of all entries in the database.
* POST:
  * /person: Creates a new entry with an ID of 1 higher than the highest ID in the database. Input is expected in JSON format.
  * /person/{id}: Creates a new entry for a specific ID. Input is expected in JSON format.
  * /import: Accepts CSV formatted data, which is imported into the database.
* PUT:
  * /person/{id}: Replaces the current entry with the provided information in JSON format.
* PATCH:
  * /person/{id}: Updates the given information for an entry. Accepts partial information.
* DELETE:
  * /person/{id}: Deletes a specified entry from the database.