package ports

import (
	"net/http"

	"github.com/go-chi/render"
)

type HttpServer struct{}

// (GET /trainer/calendar)
func (_ HttpServer) GetTrainerAvailableHours(w http.ResponseWriter, r *http.Request, params GetTrainerAvailableHoursParams) {
	w.WriteHeader(http.StatusOK)
	render.Respond(w, r, map[string]string{
		"message": "This endpoint is not implemented yet",
	})
}

// (PUT /trainer/calendar/make-hour-available)
func (_ HttpServer) MakeHourAvailable(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (PUT /trainer/calendar/make-hour-unavailable)
func (_ HttpServer) MakeHourUnavailable(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	render.Respond(w, r, map[string]string{
		"message": "This endpoint is not implemented yet",
	})
}
