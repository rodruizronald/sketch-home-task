package router

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Router struct {
	http.Handler
	muxRouter    *mux.Router
	validator    *validator.Validate
	templatesDir string
}

type HandlerRequest struct {
	Context context.Context
	Vars    map[string]string
	Body    interface{}
}

type HandlerResponse struct {
	Response    interface{}
	ContentType ContentType
	Status      int
	// Name of the template file to execute
	Template string
}

func (h *HandlerResponse) SetText(resp string, status int) {
	h.Response = resp
	h.ContentType = ContentTypeText
	h.Status = status
}

func (h *HandlerResponse) SetHTML(resp interface{}, template string, status int) {
	h.Response = resp
	h.Template = template
	h.ContentType = ContentTypeHTML
	h.Status = status
}

type HandlerFunc func(req *HandlerRequest) (resp *HandlerResponse)

type ContentType int

const (
	ContentTypeText ContentType = iota
	ContentTypeHTML
	ContentTypeJSON
)

func NewRouter(validator *validator.Validate, templatesDir string) (r *Router) {
	muxRouter := mux.NewRouter()
	muxSubrouter := muxRouter.Methods(
		http.MethodPost,
		http.MethodPut,
		http.MethodGet,
		http.MethodDelete).Subrouter()

	return &Router{
		Handler:      muxRouter,
		muxRouter:    muxSubrouter,
		validator:    validator,
		templatesDir: templatesDir,
	}
}

// ----------------------- Request Handler Methods ----------------------- //

func (r *Router) POST(path string, body interface{}, handler HandlerFunc) {
	if reflect.ValueOf(body).Kind() != reflect.Ptr {
		panic("body is not addressable")
	}
	r.handle(http.MethodPost, path, body, handler)
}

func (r *Router) PUT(path string, body interface{}, handler HandlerFunc) {
	if reflect.ValueOf(body).Kind() != reflect.Ptr {
		panic("body is not addressable")
	}
	r.handle(http.MethodPut, path, body, handler)
}

func (r *Router) GET(path string, handler HandlerFunc) {
	r.handle(http.MethodGet, path, nil, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.handle(http.MethodDelete, path, nil, handler)
}

// ----------------------- Router Request Handler ----------------------- //

func (r *Router) handle(method string, path string, body interface{}, handler HandlerFunc) {
	r.muxRouter.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if http.NoBody == req.Body && (method == http.MethodPost || method == http.MethodPut) {
			writeError(w, nil, "request content is empty", http.StatusBadRequest)
			return
		}

		// Only JSON requests supported
		if http.NoBody != req.Body {
			defer req.Body.Close()

			rawBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				writeError(w, err, "unable to read request content", http.StatusInternalServerError)
				return
			}
			if err = json.Unmarshal(rawBody, body); err != nil {
				writeError(w, err, "failed to process request content", http.StatusBadRequest)
				return
			}
			if r.validator != nil {
				if err = r.validator.Struct(body); err != nil {
					writeError(w, err, "request content is invalid", http.StatusBadRequest)
					return
				}
			}
		}

		handlerReq := &HandlerRequest{
			Context: req.Context(),
			Vars:    mux.Vars(req),
			Body:    body,
		}

		handler(handlerReq).writeResponse(w, r.templatesDir)
	}).Methods(method)
}

// ----------------------- Router Response Writer ----------------------- //

func (h *HandlerResponse) writeResponse(w http.ResponseWriter, templatesDir string) {
	var resp string
	var contentType string

	switch h.ContentType {
	case ContentTypeText:
		resp = h.Response.(string)
		contentType = "text/plain; charset=utf-8"
	case ContentTypeJSON:
		respBytes, err := json.MarshalIndent(h.Response, "", "  ")
		if err != nil {
			writeError(w, err, "failed to process respond content", http.StatusInternalServerError)
			return
		}
		resp = string(respBytes)
		contentType = "application/json; charset=utf-8"
	case ContentTypeHTML:
		templatePath := filepath.Join(templatesDir, h.Template)
		tpl, err := template.New(h.Template).ParseFiles(templatePath)
		if err != nil {
			writeError(w, err, "failed to parse template files", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(h.Status)
		err = tpl.Execute(w, h.Response)
		if err != nil {
			writeError(w, err, "failed to execute template", http.StatusInternalServerError)
			return
		}
		return
	default:
		panic("unknown content-type")
	}

	contentLength := strconv.FormatInt(int64(len(resp)), 10)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", contentLength)
	w.WriteHeader(h.Status)
	w.Write([]byte(resp))
}

func writeError(w http.ResponseWriter, internalErr error, httpErr string, status int) {
	if internalErr != nil {
		log.Printf("[ERROR] %v: %v\n", httpErr, internalErr)
	}
	http.Error(w, httpErr, status)
}
