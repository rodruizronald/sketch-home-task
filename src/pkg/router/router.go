package router

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Router struct {
	http.Handler
	muxRouter *mux.Router
	validator *validator.Validate
	logError  bool
}

type HandlerFunc func(vars map[string]string, body interface{}) (resp interface{}, status int)

const (
	MethodPost   string = "POST"
	MethodPut    string = "PUT"
	MethodGet    string = "GET"
	MethodDelete string = "DELETE"
)

func NewRouter(validator *validator.Validate, logError bool) (r *Router) {
	muxRouter := mux.NewRouter()
	muxSubrouter := muxRouter.Methods(
		MethodPost,
		MethodPut,
		MethodGet,
		MethodDelete).Subrouter()

	return &Router{
		Handler:   muxRouter,
		muxRouter: muxSubrouter,
		validator: validator,
		logError:  logError,
	}
}

func (r *Router) POST(path string, body interface{}, handler HandlerFunc) {
	if reflect.ValueOf(body).Kind() != reflect.Ptr {
		panic("body is not addressable")
	}
	r.handle(MethodPost, path, body, handler)
}

func (r *Router) PUT(path string, body interface{}, handler HandlerFunc) {
	if reflect.ValueOf(body).Kind() != reflect.Ptr {
		panic("body is not addressable")
	}
	r.handle(MethodPut, path, body, handler)
}

func (r *Router) GET(path string, handler HandlerFunc) {
	r.handle(MethodGet, path, nil, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.handle(MethodDelete, path, nil, handler)
}

func (r *Router) handle(method string, path string, body interface{}, handler HandlerFunc) {
	r.muxRouter.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if http.NoBody == req.Body && (method == MethodPost || method == MethodPut) {
			r.writeError(w, nil, "request content is empty", http.StatusBadRequest)
			return
		}

		if http.NoBody != req.Body {
			defer req.Body.Close()

			rawBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				r.writeError(w, err, "unable to read request content", http.StatusInternalServerError)
				return
			}
			if err = json.Unmarshal(rawBody, body); err != nil {
				r.writeError(w, err, "failed to process request content", http.StatusBadRequest)
				return
			}
			if r.validator != nil {
				if err = r.validator.Struct(body); err != nil {
					r.writeError(w, err, "request content is invalid", http.StatusBadRequest)
					return
				}
			}
		} else if method == MethodPost || method == MethodPut {
			r.writeError(w, nil, "request content is empty", http.StatusBadRequest)
			return
		}

		response, status := handler(mux.Vars(req), body)

		var contentType string
		var responseStr string

		switch v := response.(type) {
		case string:
			responseStr = v
			contentType = "text/plain; charset=utf-8"
		default:
			responseBytes, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				r.writeError(w, err, "failed to process respond content", http.StatusInternalServerError)
				return
			}
			responseStr = string(responseBytes)
			contentType = "application/json; charset=utf-8"
		}

		contentLength := strconv.FormatInt(int64(len(responseStr)), 10)

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", contentLength)
		w.WriteHeader(status)
		w.Write([]byte(responseStr))
	}).Methods(method)
}

func (r *Router) writeError(w http.ResponseWriter, internalErr error, httpErr string, status int) {
	if r.logError && internalErr != nil {
		log.Printf("[ERROR] %v: %v\n", httpErr, internalErr)
	}
	http.Error(w, httpErr, status)
}
