package ports

import (
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/app/command"
	"net/http"

	"github.com/go-chi/render"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/quangbach27/wild-workout-microservice/internal/common/server/httperr"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/app"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/app/query"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

func (httpServer HttpServer) GetTrainerAvailableHours(w http.ResponseWriter, r *http.Request, params GetTrainerAvailableHoursParams) {
	dateApps, err := httpServer.app.Queries.AvailableHours.Handle(r.Context(), query.AvailableHours{
		From: params.DateFrom,
		To:   params.DateTo,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	dates := dateModelsToResponse(dateApps)
	render.Respond(w, r, dates)
}

// (PUT /trainer/calendar/make-hour-available)
func (httpServer HttpServer) MakeHourAvailable(w http.ResponseWriter, r *http.Request) {
	hourUpdate := &HourUpdate{}
	if err := render.Decode(r, hourUpdate); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err := httpServer.app.Commands.MakeHoursAvailable.Handle(r.Context(), command.MakeHoursAvailableCommand{Hours: hourUpdate.Hours})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// (PUT /trainer/calendar/make-hour-unavailable)
func (httpSever HttpServer) MakeHourUnavailable(w http.ResponseWriter, r *http.Request) {
	hourUpdate := &HourUpdate{}
	if err := render.Decode(r, hourUpdate); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err := httpSever.app.Commands.MakeHoursNotAvailable.Handle(r.Context(), command.MakeHoursNotAvailableCommand{Hours: hourUpdate.Hours})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
	}

	w.WriteHeader(http.StatusNoContent)
}

func dateModelsToResponse(models []query.Date) []Date {
	var dates []Date
	for _, d := range models {
		var hours []Hour
		for _, h := range d.Hours {
			hours = append(hours, Hour{
				Available:            h.Available,
				HasTrainingScheduled: h.HasTrainingScheduled,
				Hour:                 h.Hour,
			})
		}

		dates = append(dates, Date{
			Date: openapi_types.Date{
				Time: d.Date,
			},
			HasFreeHours: d.HasFreeHours,
			Hours:        hours,
		})
	}

	return dates
}
