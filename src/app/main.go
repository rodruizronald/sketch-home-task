package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator"
	_ "github.com/lib/pq"
	"github.com/sketch-home-task/src/pkg/dba"
	"github.com/sketch-home-task/src/pkg/illustrator"
	"github.com/sketch-home-task/src/pkg/router"
)

type App struct {
	router  *router.Router
	storage illustrator.CanvasStorage
}

func main() {
	serverPort := os.Getenv("SERVER_PORT")
	templatesDir := os.Getenv("TEMPLATES_DIRECTORY")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDatabase := os.Getenv("POSTGRES_DATABASE")

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresHost, postgresPort, postgresUser, postgresPassword, postgresDatabase)

	storage, err := dba.NewStorage("postgres", dns, 1, 1)
	if err != nil {
		panic(err)
	}

	validator := validator.New()
	illustrator.RegisterValidation(validator)
	router := router.NewRouter(validator, templatesDir)

	app := App{
		router:  router,
		storage: storage,
	}

	// Register canvas API end points
	app.router.POST("/canvas", &illustrator.CanvasModel{}, app.createCanvas)
	app.router.PUT("/canvas", &illustrator.CanvasModel{}, app.updateCanvas)
	app.router.GET("/canvas/{name:[a-z]{1,25}}", app.getCanvas)
	app.router.DELETE("/canvas/{name:[a-z]{1,25}}", app.deleteCanvas)

	addr := fmt.Sprintf(":%s", serverPort)
	srv := http.Server{
		Addr:    addr,
		Handler: app.router,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}

func (a *App) createCanvas(req *router.HandlerRequest) (resp *router.HandlerResponse) {
	resp = new(router.HandlerResponse)
	canvas := getCanvasFromRequest(req)

	_, err := a.storage.Create(req.Context, canvas)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			resp.SetText("canvas name already exists", http.StatusBadRequest)
		} else {
			setInternalErrorResponse(resp, "failed to create canvas", err)
		}
		return
	}

	resp.SetText("create OK", http.StatusCreated)
	return
}

func (a *App) updateCanvas(req *router.HandlerRequest) (resp *router.HandlerResponse) {
	resp = new(router.HandlerResponse)
	canvas := getCanvasFromRequest(req)

	res, err := a.storage.Update(req.Context, canvas)
	if err != nil {
		setInternalErrorResponse(resp, "failed to update canvas", err)
		return
	}

	count, err := res.RowsAffected()
	if err == nil && count > 0 {
		resp.SetText("update OK", http.StatusOK)
		return
	}

	resp.SetText("canvas not found", http.StatusBadRequest)
	return
}

func (a *App) getCanvas(req *router.HandlerRequest) (resp *router.HandlerResponse) {
	resp = new(router.HandlerResponse)

	name, ok := req.Vars["name"]
	if !ok {
		resp.SetText("route variable 'name' not found", http.StatusBadRequest)
		return
	}

	canvas, err := a.storage.FindByName(req.Context, hashString(name))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp.SetText("canvas not found", http.StatusBadRequest)
		} else {
			setInternalErrorResponse(resp, "failed to retrieve canvas", err)
		}
		return
	}

	// Validation is carried out in the router
	str, _ := canvas.GetString(' ', "<br>", nil)
	templateData := &struct{ Canvas template.HTML }{template.HTML(str)}
	resp.SetHTML(templateData, "index.html", http.StatusOK)
	return
}

func (a *App) deleteCanvas(req *router.HandlerRequest) (resp *router.HandlerResponse) {
	resp = new(router.HandlerResponse)

	name, ok := req.Vars["name"]
	if !ok {
		resp.SetText("route variable 'name' not found", http.StatusBadRequest)
		return
	}

	res, err := a.storage.Delete(req.Context, hashString(name))
	if err != nil {
		setInternalErrorResponse(resp, "failed to delete canvas", err)
		return
	}

	count, err := res.RowsAffected()
	if err == nil && count > 0 {
		resp.SetText("delete OK", http.StatusOK)
		return
	}

	resp.SetText("canvas not found", http.StatusBadRequest)
	return
}

func getCanvasFromRequest(req *router.HandlerRequest) (canvas *illustrator.CanvasModel) {
	canvas = req.Body.(*illustrator.CanvasModel)
	canvas.Name = hashString(canvas.Name)
	return
}

func setInternalErrorResponse(resp *router.HandlerResponse, msg string, err error) {
	resp.SetText(msg, http.StatusInternalServerError)
	log.Printf("[ERROR] %v: %v\n", msg, err)
}

func hashString(str string) (result string) {
	hash := sha1.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}
