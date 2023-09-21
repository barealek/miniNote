package controllers

import (
	"github.com/barealek/mininote/internal/controllers/root"
	"github.com/go-chi/chi/v5"
)

func RegisterControllers(app *chi.Mux) {
	app.Route("/", root.Router)
}
