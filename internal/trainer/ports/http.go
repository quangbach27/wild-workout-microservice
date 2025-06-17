package ports

import (
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

func (http HttpServer) GetTrainerAvailableHours(w http.ResponseWriter, r *http.Request, params GetTrainerAvailableHoursParams) {
	dateApps, err := http.app.Queries.AvailableHours.Handle(r.Context(), query.AvailableHours{
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
func (_ HttpServer) MakeHourAvailable(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (PUT /trainer/calendar/make-hour-unavailable)
func (_ HttpServer) MakeHourUnavailable(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
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
