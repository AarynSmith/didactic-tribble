package app

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func executeRequest(req *http.Request, create bool) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a := App{}
	err := a.Initialize("Test.sqlitedb")
	clearTable(a.Database)
	if create {
		p := Person{
			ID:        1,
			FirstName: "Test",
			LastName:  "Name",
			Email:     "Test.Name@example.com",
			Phone:     "123-456-7890",
		}
		p.dbCreatePerson(a.Database)
	}
	if err != nil {
		log.Fatalf("Error Initializing: %v", err.Error())
	}
	a.addHandles()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func TestIndex(t *testing.T) {
	tests := []struct {
		request      string
		expectedCode int
	}{
		{
			request:      "/",
			expectedCode: 200,
		},
		{
			request:      "/NonExistantHandle",
			expectedCode: 404,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.request, nil)
		response := executeRequest(req, true)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestNotAllowed(t *testing.T) {
	tests := []struct {
		request      string
		method       string
		expectedCode int
	}{
		{
			request:      "/people",
			method:       "GET",
			expectedCode: 200,
		},
		{
			request:      "/people",
			method:       "POST",
			expectedCode: 405,
		},
		{
			request:      "/people",
			method:       "DELETE",
			expectedCode: 405,
		},
		{
			request:      "/people",
			method:       "PUT",
			expectedCode: 405,
		},
		{
			request:      "/people",
			method:       "UPDATE",
			expectedCode: 405,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest(tt.method, tt.request, nil)
		response := executeRequest(req, true)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_CreatePersonForm(t *testing.T) {
	tests := []struct {
		request      string
		expectedCode int
	}{
		{
			request:      "/people/create",
			expectedCode: 200,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.request, nil)
		response := executeRequest(req, true)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

// func TestApp_UpdatePersonForm(t *testing.T) {
// 	tests := []struct {
// 		request      string
// 		expectedCode int
// 		nonempty     bool
// 	}{
// 		{
// 			request:      "/people/1/update",
// 			expectedCode: 200,
// 			nonempty:     true,
// 		},
// 		{
// 			request:      "/people/1/update/",
// 			expectedCode: 404,
// 			nonempty:     false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		req, _ := http.NewRequest("GET", tt.request, nil)
// 		response := executeRequest(req, tt.nonempty)
// 		if tt.expectedCode != response.Code {
// 			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
// 		}
// 	}
// }

func TestApp_GetPeople(t *testing.T) {
	tests := []struct {
		request      string
		expectedCode int
		nonempty     bool
	}{
		{
			request:      "/people",
			expectedCode: 200,
			nonempty:     true,
		},
		{
			request:      "/people",
			expectedCode: 200,
			nonempty:     false,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.request, nil)
		response := executeRequest(req, true)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_GetPerson(t *testing.T) {
	tests := []struct {
		request      string
		expectedCode int
		nonempty     bool
	}{
		{
			request:      "/people/1",
			expectedCode: 200,
			nonempty:     true,
		},
		{
			request:      "/people/1",
			expectedCode: 404,
			nonempty:     false,
		},
		{
			request:      "/people/99",
			expectedCode: 404,
			nonempty:     true,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.request, nil)
		response := executeRequest(req, tt.nonempty)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}
func TestApp_CreatePerson(t *testing.T) {
	tests := []struct {
		request      string
		form         url.Values
		expectedCode int
		nonempty     bool
	}{
		{
			request: "/people/1",
			form: url.Values{
				"id":    []string{"1"},
				"fname": []string{"Test"},
				"lname": []string{"User"},
				"email": []string{"TestUser@example.com"},
				"phone": []string{"987-654-3210"},
			},
			expectedCode: 409,
			nonempty:     true,
		},
		{
			request: "/people/1",
			form: url.Values{
				"id":    []string{"1"},
				"fname": []string{"Test"},
				"lname": []string{"User"},
				"email": []string{"TestUser@example.com"},
				"phone": []string{"987-654-3210"},
			},
			expectedCode: 200,
			nonempty:     false,
		},
		{
			request: "/people/2",
			form: url.Values{
				"id":    []string{"2"},
				"fname": []string{"Test"},
				"lname": []string{"User 2"},
				"email": []string{"TestUser2@example.com"},
				"phone": []string{"987-654-3210"},
			},
			expectedCode: 200,
			nonempty:     true,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("POST", tt.request, nil)
		req.PostForm = tt.form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response := executeRequest(req, tt.nonempty)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_DeletePerson(t *testing.T) {
	tests := []struct {
		request      string
		expectedCode int
		nonempty     bool
	}{
		{
			request:      "/people/1",
			expectedCode: 200,
			nonempty:     true,
		},
		{
			request:      "/people/1",
			expectedCode: 200,
			nonempty:     false,
		},
		{
			request:      "/people/99",
			expectedCode: 200,
			nonempty:     true,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("DELETE", tt.request, nil)
		response := executeRequest(req, tt.nonempty)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_Initialize(t *testing.T) {
	type fields struct {
		Router   *mux.Router
		Database *sql.DB
	}
	type args struct {
		dbname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Database",
			args: args{
				dbname: "test.sqlitedb",
			},
		},
		{
			name: "Null Database",
			args: args{
				dbname: ".",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				Router:   tt.fields.Router,
				Database: tt.fields.Database,
			}
			if err := a.Initialize(tt.args.dbname); (err != nil) != tt.wantErr {
				t.Errorf("App.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
			os.Remove(tt.args.dbname)
		})
	}
}
