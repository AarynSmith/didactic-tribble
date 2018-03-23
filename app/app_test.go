package app

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

const TestDBName = "Test.sqlitedb"

func executeRequest(req *http.Request, emptydb bool) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a := App{}
	err := a.Initialize(TestDBName)
	clearTable(a.Database)
	if !emptydb {
		p := Person{
			id:        1,
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
	os.Remove(TestDBName)
	return rr
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

func TestApp_ReadPeople(t *testing.T) {
	tests := []struct {
		request      string
		method       string
		expectedCode int
		emptydb      bool
	}{
		{
			request:      "/people",
			method:       "GET",
			expectedCode: 200,
		},
		{
			request:      "/people",
			method:       "GET",
			expectedCode: 500,
			emptydb:      true,
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
		log.Printf("%v %v %v", tt.method, tt.request, tt.emptydb)
		req, _ := http.NewRequest(tt.method, tt.request, nil)
		response := executeRequest(req, tt.emptydb)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_CreatePerson(t *testing.T) {
	tests := []struct {
		request      string
		body         string
		expectedCode int
		emptydb      bool
	}{
		{
			request: "/person/1",
			body: `{
				"FirstName": "Test",
				"LastName": "User",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 409,
		},
		{
			request: "/person/1",
			body: `{
				"FirstName": "Test",
				"LastName": "User",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 200,
			emptydb:      true,
		},
		{
			request: "/person/2",
			body: `{
				"FirstName": "Test",
				"LastName": "User 2",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 200,
		},
		{
			request:      "/person/2",
			body:         `BadJson`,
			expectedCode: 500,
		},
		{
			request: "/person",
			body: `{
				"FirstName": "Test",
				"LastName": "User 2",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 200,
		},
		{
			request: "/person/9223372036854775808",
			body: `{
				"FirstName": "Test",
				"LastName": "User 2",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 409,
		},
	}
	for _, tt := range tests {
		buf := strings.NewReader(tt.body)
		req, _ := http.NewRequest("POST", tt.request, buf)
		response := executeRequest(req, tt.emptydb)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_ReadPerson(t *testing.T) {
	tests := []struct {
		request      string
		expectedCode int
		emptydb      bool
	}{
		{
			request:      "/person/1",
			expectedCode: 200,
		},
		{
			request:      "/person/1",
			expectedCode: 404,
			emptydb:      true,
		},
		{
			request:      "/person/99",
			expectedCode: 404,
		},
		{
			request:      "/person/9223372036854775809",
			expectedCode: 409,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.request, nil)
		response := executeRequest(req, tt.emptydb)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_UpdatePerson(t *testing.T) {
	tests := []struct {
		request      string
		body         string
		expectedCode int
		emptydb      bool
	}{
		{
			request: "/person/1",
			body: `{
				"FirstName": "Test",
				"LastName": "User",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 200,
		},
		{
			request: "/person/1",
			body: `{
				"FirstName": "Test",
				"LastName": "User",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 404,
			emptydb:      true,
		},
		{
			request: "/person/2",
			body: `{
				"FirstName": "Test",
				"LastName": "User 2",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 404,
		},
		{
			request:      "/person/2",
			body:         `BadJson`,
			expectedCode: 500,
		},
		{
			request:      "/person/9223372036854775809",
			body:         ``,
			expectedCode: 409,
		},
	}
	for _, tt := range tests {
		buf := strings.NewReader(tt.body)
		req, _ := http.NewRequest("PUT", tt.request, buf)
		response := executeRequest(req, tt.emptydb)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_UpdatePatchPerson(t *testing.T) {
	tests := []struct {
		request      string
		body         string
		expectedCode int
		emptydb      bool
	}{
		{
			request: "/person/1",
			body: `{
				"FirstName": "Test",
				"LastName": "User",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 200,
		},
		{
			request: "/person/1",
			body: `{
				"FirstName": "Test",
				"LastName": "User",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 404,
			emptydb:      true,
		},
		{
			request: "/person/2",
			body: `{
				"FirstName": "Test",
				"LastName": "User 2",
				"Email": "TestUser@example.com",
				"Phone": "987-654-3210"
			}`,
			expectedCode: 404,
		},
		{
			request:      "/person/2",
			body:         `BadJson`,
			expectedCode: 500,
		},
		{
			request:      "/person/9223372036854775809",
			body:         ``,
			expectedCode: 409,
		},
	}
	for _, tt := range tests {
		buf := strings.NewReader(tt.body)
		req, _ := http.NewRequest("PATCH", tt.request, buf)
		response := executeRequest(req, tt.emptydb)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_DeletePerson(t *testing.T) {
	tests := []struct {
		request      string
		expectedCode int
		emptydb      bool
	}{
		{
			request:      "/person/1",
			expectedCode: 200,
		},
		{
			request:      "/person/1",
			expectedCode: 200,
			emptydb:      true,
		},
		{
			request:      "/person/99",
			expectedCode: 200,
		},
		{
			request:      "/person/9223372036854775809",
			expectedCode: 409,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("DELETE", tt.request, nil)
		response := executeRequest(req, tt.emptydb)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}
func TestApp_Import(t *testing.T) {
	tests := []struct {
		method       string
		request      string
		body         string
		expectedCode int
		emptydb      bool
	}{
		{
			method:  "POST",
			request: "/import",
			body: `FirstName,LastName,Email,Phone
Person 1,Smith,email@example.com,123-456-7890
Person 2,Smith,email@example.com,123-456-7890
Person 3,Smith,email@example.com,123-456-7890
Person 4,Smith,email@example.com,123-456-7890
`,
			expectedCode: 200,
		},
		{
			method:  "POST",
			request: "/import",
			body: `FirstName,LastName,Email,Phone
Person 1,Smith,email@example.com,123-456-7890
Person 2,Smith,email@example.com,123-456-7890
Person 3,Smith,email@example.com,123-456-7890
Person 4,Smith,email@example.com,123-456-7890
`,
			expectedCode: 200,
			emptydb:      true,
		},
		{
			method:       "POST",
			request:      "/import",
			body:         "",
			expectedCode: 409,
		},
	}
	for _, tt := range tests {
		buf := strings.NewReader(tt.body)
		req, _ := http.NewRequest(tt.method, tt.request, buf)
		response := executeRequest(req, tt.emptydb)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}

func TestApp_Export(t *testing.T) {
	tests := []struct {
		method       string
		request      string
		expectedCode int
		emptydb      bool
	}{
		{
			method:       "GET",
			request:      "/export",
			expectedCode: 200,
		},
		{
			method:       "GET",
			request:      "/export",
			expectedCode: 500,
			emptydb:      true,
		},
		{
			method:       "POST",
			request:      "/export",
			expectedCode: 405,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest(tt.method, tt.request, nil)
		response := executeRequest(req, tt.emptydb)
		if tt.expectedCode != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", tt.expectedCode, response.Code)
		}
	}
}
