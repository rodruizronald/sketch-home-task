package router_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/sketch-demo/src/pkg/router"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Data   string `json:"data" validate:"alpha"`
	Length int    `json:"length" validate:"gt=5"`
}

type testEntry struct {
	name        string
	status      int
	uri         string
	method      string
	reqContent  interface{}
	respContent interface{}
}

func TestRouterPostHandler(t *testing.T) {
	testPath := "/data"
	validReqContent := TestData{
		Data:   "gooddata",
		Length: 8,
	}

	testTable := []testEntry{
		// Valid test cases
		{
			name:        "Test with valid post req",
			status:      http.StatusCreated,
			method:      router.MethodPost,
			uri:         "/data",
			reqContent:  validReqContent,
			respContent: nil,
		},
		// Invalid test cases
		{
			name:   "Test with unmatching req content",
			status: http.StatusBadRequest,
			uri:    "/data",
			method: router.MethodPost,
			reqContent: TestData{
				Data:   "baddata",
				Length: 7,
			},
			respContent: nil,
		},
		{
			name:   "Test with invalid req content",
			status: http.StatusBadRequest,
			uri:    "/data",
			method: router.MethodPost,
			reqContent: TestData{
				Data:   "112!@$#%ffwr",
				Length: 3,
			},
			respContent: nil,
		},
		{
			name:        "Test with invalid method",
			status:      http.StatusMethodNotAllowed,
			uri:         "/data",
			method:      router.MethodPut,
			reqContent:  validReqContent,
			respContent: nil,
		},
		{
			name:        "Test with invalid uri",
			status:      http.StatusNotFound,
			uri:         "/invalid",
			method:      router.MethodPut,
			reqContent:  validReqContent,
			respContent: nil,
		},
		{
			name:        "Test with empty body",
			status:      http.StatusBadRequest,
			uri:         "/data",
			method:      router.MethodPost,
			reqContent:  nil,
			respContent: nil,
		},
	}

	validator := validator.New()
	handler := router.NewRouter(validator, true)

	// Validate data correctness in POST handler
	handler.POST(testPath, &TestData{}, func(vars map[string]string, body interface{}) (resp interface{}, status int) {
		if reflect.DeepEqual(validReqContent, reflect.ValueOf(body).Elem().Interface()) {
			return "ok content", http.StatusCreated
		}
		return "bad content", http.StatusBadRequest
	})

	runRouterTests(t, testTable, handler, testPath)
}

func TestRouterPutHandler(t *testing.T) {
	testPath := "/data/{data_id:[0-9]+}"
	validReqContent := TestData{
		Data:   "gooddata",
		Length: 8,
	}

	testTable := []testEntry{
		// Valid test cases
		{
			name:        "Test with valid put req",
			status:      http.StatusOK,
			method:      router.MethodPut,
			uri:         "/data/25",
			reqContent:  validReqContent,
			respContent: nil,
		},
		// Invalid test cases
		{
			name:   "Test with unmatching req content",
			status: http.StatusBadRequest,
			uri:    "/data/25",
			method: router.MethodPut,
			reqContent: TestData{
				Data:   "baddata",
				Length: 7,
			},
			respContent: nil,
		},
		{
			name:   "Test with invalid req content",
			status: http.StatusBadRequest,
			uri:    "/data/25",
			method: router.MethodPut,
			reqContent: TestData{
				Data:   "112!@$#%ffwr",
				Length: 3,
			},
			respContent: nil,
		},
		{
			name:        "Test with invalid method",
			status:      http.StatusMethodNotAllowed,
			uri:         "/data/25",
			method:      router.MethodDelete,
			reqContent:  validReqContent,
			respContent: nil,
		},
		{
			name:        "Test with invalid uri",
			status:      http.StatusNotFound,
			uri:         "/data/2a2#",
			method:      router.MethodPut,
			reqContent:  validReqContent,
			respContent: nil,
		},
		{
			name:        "Test with empty body",
			status:      http.StatusBadRequest,
			uri:         "/data/25",
			method:      router.MethodPut,
			reqContent:  nil,
			respContent: nil,
		},
	}

	validator := validator.New()
	handler := router.NewRouter(validator, true)

	// Validate data correctness in PUT handler
	handler.PUT(testPath, &TestData{}, func(vars map[string]string, body interface{}) (resp interface{}, status int) {
		if reflect.DeepEqual(validReqContent, reflect.ValueOf(body).Elem().Interface()) {
			return "ok content", http.StatusOK
		}
		return "bad content", http.StatusBadRequest
	})

	runRouterTests(t, testTable, handler, testPath)
}

func TestRouterGetHandler(t *testing.T) {
	testPath := "/data/{data_id:[0-9]+}"
	validRespContent := TestData{
		Data:   "gooddata",
		Length: 8,
	}

	testTable := []testEntry{
		// Valid test cases
		{
			name:        "Test with valid get req",
			status:      http.StatusOK,
			method:      router.MethodGet,
			uri:         "/data/25",
			reqContent:  nil,
			respContent: validRespContent,
		},
		// Invalid test cases
		{
			name:        "Test with invalid method",
			status:      http.StatusMethodNotAllowed,
			uri:         "/data/25",
			method:      router.MethodPost,
			reqContent:  nil,
			respContent: validRespContent,
		},
		{
			name:        "Test with invalid uri",
			status:      http.StatusNotFound,
			uri:         "/data/2a2#",
			method:      router.MethodGet,
			reqContent:  nil,
			respContent: validRespContent,
		},
	}

	validator := validator.New()
	handler := router.NewRouter(validator, true)

	// Validate data correctness in GET handler
	handler.GET(testPath, func(vars map[string]string, body interface{}) (resp interface{}, status int) {
		value, ok := vars["data_id"]
		if ok && value == "25" {
			return validRespContent, http.StatusOK
		}
		return "bad vars", http.StatusBadRequest
	})

	runRouterTests(t, testTable, handler, testPath)
}

func TestRouterDeleteHandler(t *testing.T) {
	testPath := "/data/{data_id:[0-9]+}"

	testTable := []testEntry{
		// Valid test cases
		{
			name:        "Test with valid delete req",
			status:      http.StatusOK,
			method:      router.MethodDelete,
			uri:         "/data/25",
			reqContent:  nil,
			respContent: nil,
		},
		// Invalid test cases
		{
			name:        "Test with invalid method",
			status:      http.StatusMethodNotAllowed,
			uri:         "/data/25",
			method:      router.MethodGet,
			reqContent:  nil,
			respContent: nil,
		},
		{
			name:        "Test with invalid uri",
			status:      http.StatusNotFound,
			uri:         "/data/2a2#",
			method:      router.MethodPut,
			reqContent:  nil,
			respContent: nil,
		},
	}

	validator := validator.New()
	handler := router.NewRouter(validator, true)

	// Validate data correctness in DELETE handler
	handler.DELETE(testPath, func(vars map[string]string, body interface{}) (resp interface{}, status int) {
		value, ok := vars["data_id"]
		if ok && value == "25" {
			return "ok vars", http.StatusOK
		}
		return "bad vars", http.StatusBadRequest
	})

	runRouterTests(t, testTable, handler, testPath)
}

func runRouterTests(t *testing.T, testTable []testEntry, handler http.Handler, testPath string) {
	a := assert.New(t)

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(handler)
			defer server.Close()

			// Add request body if any
			reader := &bytes.Reader{}
			if tt.reqContent != nil {
				b, err := json.MarshalIndent(&tt.reqContent, "", "")
				a.NoError(err)
				reader = bytes.NewReader(b)
			}

			request, err := http.NewRequest(tt.method, server.URL+tt.uri, reader)
			a.NoError(err)

			resp, err := server.Client().Do(request)
			a.NoError(err)

			// Validate status code
			a.Equal(tt.status, resp.StatusCode)

			// Validate body content
			if http.NoBody != resp.Body && tt.respContent != nil {
				defer resp.Body.Close()

				actualRawBody, err := io.ReadAll(resp.Body)
				a.NoError(err)

				expectedRawBody, err := json.MarshalIndent(tt.respContent, "", "")
				a.NoError(err)

				expected := strings.TrimSpace(string(expectedRawBody))
				actual := strings.TrimSpace(string(actualRawBody))

				a.NotEqual(expected, actual)
			}
		})
	}
}
