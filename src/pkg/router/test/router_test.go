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
	"github.com/sketch-home-task/src/pkg/router"
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
	contentType router.ContentType
}

func TestRouterPostHandler(t *testing.T) {
	testPath := "/data"
	validJSONReq := TestData{
		Data:   "gooddata",
		Length: 8,
	}

	testTable := []testEntry{
		// Valid test cases
		{
			name:       "Test with valid post req",
			status:     http.StatusCreated,
			method:     http.MethodPost,
			uri:        "/data",
			reqContent: validJSONReq,
		},
		// Invalid test cases
		{
			name:   "Test with unmatching req content",
			status: http.StatusBadRequest,
			uri:    "/data",
			method: http.MethodPost,
			reqContent: TestData{
				Data:   "baddata",
				Length: 7,
			},
		},
		{
			name:   "Test with invalid req content",
			status: http.StatusBadRequest,
			uri:    "/data",
			method: http.MethodPost,
			reqContent: TestData{
				Data:   "112!@$#%ffwr",
				Length: 3,
			},
		},
		{
			name:       "Test with invalid method",
			status:     http.StatusMethodNotAllowed,
			uri:        "/data",
			method:     http.MethodPut,
			reqContent: validJSONReq,
		},
		{
			name:       "Test with invalid uri",
			status:     http.StatusNotFound,
			uri:        "/invalid",
			method:     http.MethodPut,
			reqContent: validJSONReq,
		},
		{
			name:   "Test with empty body",
			status: http.StatusBadRequest,
			uri:    "/data",
			method: http.MethodPost,
		},
	}

	validator := validator.New()
	handler := router.NewRouter(validator)

	// Validate data correctness in POST handler
	handler.POST(testPath, &TestData{}, func(req *router.HandlerRequest) (resp *router.HandlerResponse) {
		resp = new(router.HandlerResponse)
		resp.ContentType = router.ContentTypeText
		resp.Response = "bad content"
		resp.Status = http.StatusBadRequest

		if reflect.DeepEqual(validJSONReq, reflect.ValueOf(req.Body).Elem().Interface()) {
			resp.Response = "ok content"
			resp.Status = http.StatusCreated
		}

		return
	})

	runRouterTests(t, testTable, handler)
}

func TestRouterPutHandler(t *testing.T) {
	testPath := "/data/{data_id:[0-9]+}"
	validJSONReq := TestData{
		Data:   "gooddata",
		Length: 8,
	}

	testTable := []testEntry{
		// Valid test cases
		{
			name:       "Test with valid put req",
			status:     http.StatusOK,
			method:     http.MethodPut,
			uri:        "/data/25",
			reqContent: validJSONReq,
		},
		// Invalid test cases
		{
			name:   "Test with unmatching req content",
			status: http.StatusBadRequest,
			uri:    "/data/25",
			method: http.MethodPut,
			reqContent: TestData{
				Data:   "baddata",
				Length: 7,
			},
		},
		{
			name:   "Test with invalid req content",
			status: http.StatusBadRequest,
			uri:    "/data/25",
			method: http.MethodPut,
			reqContent: TestData{
				Data:   "112!@$#%ffwr",
				Length: 3,
			},
		},
		{
			name:       "Test with invalid method",
			status:     http.StatusMethodNotAllowed,
			uri:        "/data/25",
			method:     http.MethodDelete,
			reqContent: validJSONReq,
		},
		{
			name:       "Test with invalid uri",
			status:     http.StatusNotFound,
			uri:        "/data/2a2#",
			method:     http.MethodPut,
			reqContent: validJSONReq,
		},
		{
			name:   "Test with empty body",
			status: http.StatusBadRequest,
			uri:    "/data/25",
			method: http.MethodPut,
		},
	}

	validator := validator.New()
	handler := router.NewRouter(validator)

	// Validate data correctness in PUT handler
	handler.PUT(testPath, &TestData{}, func(req *router.HandlerRequest) (resp *router.HandlerResponse) {
		resp = new(router.HandlerResponse)
		resp.ContentType = router.ContentTypeText
		resp.Response = "bad content"
		resp.Status = http.StatusBadRequest

		if reflect.DeepEqual(validJSONReq, reflect.ValueOf(req.Body).Elem().Interface()) {
			resp.Response = "ok content"
			resp.Status = http.StatusOK
		}

		return
	})

	runRouterTests(t, testTable, handler)
}

func TestRouterGetHandler(t *testing.T) {
	testVarsPath := "/data/{data_id:[0-9]+}"
	testHtmlPath := "/index"

	validJSONResp := TestData{
		Data:   "testdata",
		Length: 8,
	}

	var validHTMLResp string = `<!DOCTYPE html><html lang="en"><head></head><body><p>test-canvas</p></body></html>`

	testTable := []testEntry{
		// Valid test cases
		{
			name:        "Test with valid get json req",
			status:      http.StatusOK,
			method:      http.MethodGet,
			uri:         "/data/25",
			respContent: validJSONResp,
			contentType: router.ContentTypeJSON,
		},
		{
			name:        "Test with valid get html req",
			status:      http.StatusOK,
			method:      http.MethodGet,
			uri:         testHtmlPath,
			respContent: validHTMLResp,
			contentType: router.ContentTypeHTML,
		},
		// Invalid test cases
		{
			name:   "Test with invalid method",
			status: http.StatusMethodNotAllowed,
			uri:    "/data/25",
			method: http.MethodPost,
		},
		{
			name:   "Test with invalid uri",
			status: http.StatusNotFound,
			uri:    "/data/2a2#",
			method: http.MethodGet,
		},
	}

	validator := validator.New()
	handler := router.NewRouter(validator)

	// Validate data correctness in GET handler
	handler.GET(testVarsPath, func(req *router.HandlerRequest) (resp *router.HandlerResponse) {
		resp = new(router.HandlerResponse)
		resp.ContentType = router.ContentTypeText
		resp.Response = "bad vars"
		resp.Status = http.StatusBadRequest

		value, ok := req.Vars["data_id"]
		if ok && value == "25" {
			resp.ContentType = router.ContentTypeJSON
			resp.Response = validJSONResp
			resp.Status = http.StatusOK
		}

		return
	})

	handler.GET(testHtmlPath, func(req *router.HandlerRequest) (resp *router.HandlerResponse) {
		resp = new(router.HandlerResponse)
		resp.ContentType = router.ContentTypeHTML
		resp.Status = http.StatusOK
		resp.Template = "testdata/index.tpl"
		resp.Response = &struct {
			Canvas string
		}{
			"test-canvas",
		}

		return
	})

	runRouterTests(t, testTable, handler)
}

func TestRouterDeleteHandler(t *testing.T) {
	testPath := "/data/{data_id:[0-9]+}"

	testTable := []testEntry{
		// Valid test cases
		{
			name:   "Test with valid delete req",
			status: http.StatusOK,
			method: http.MethodDelete,
			uri:    "/data/25",
		},
		// Invalid test cases
		{
			name:   "Test with invalid method",
			status: http.StatusMethodNotAllowed,
			uri:    "/data/25",
			method: http.MethodGet,
		},
		{
			name:   "Test with invalid uri",
			status: http.StatusNotFound,
			uri:    "/data/2a2#",
			method: http.MethodPut,
		},
	}

	validator := validator.New()
	handler := router.NewRouter(validator)

	// Validate data correctness in DELETE handler
	handler.DELETE(testPath, func(req *router.HandlerRequest) (resp *router.HandlerResponse) {
		resp = new(router.HandlerResponse)
		resp.ContentType = router.ContentTypeText
		resp.Response = "bad vars"
		resp.Status = http.StatusBadRequest

		value, ok := req.Vars["data_id"]
		if ok && value == "25" {
			resp.ContentType = router.ContentTypeJSON
			resp.Status = http.StatusOK
		}

		return
	})

	runRouterTests(t, testTable, handler)
}

func runRouterTests(t *testing.T, testTable []testEntry, handler http.Handler) {
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
				actualStrBody := string(actualRawBody)

				var expectedStrBody string
				switch tt.contentType {
				case router.ContentTypeText:
					expectedStrBody = tt.respContent.(string)
				case router.ContentTypeHTML:
					expectedStrBody = tt.respContent.(string)
				case router.ContentTypeJSON:
					expectedRawBody, err := json.MarshalIndent(tt.respContent, "", "")
					expectedStrBody = string(expectedRawBody)
					a.NoError(err)
				default:
					t.Error("unknown content-type")
				}

				a.Equal(strings.ReplaceAll(expectedStrBody, " ", ""), strings.ReplaceAll(actualStrBody, " ", ""))
			}
		})
	}
}
