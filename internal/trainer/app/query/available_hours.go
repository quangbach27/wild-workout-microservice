package query

import (
	"context"
	"time"

	"github.com/quangbach27/wild-workout-microservice/internal/common/cqrs"
	"github.com/quangbach27/wild-workout-microservice/internal/common/errors"
	"github.com/sirupsen/logrus"
)

type AvailableHours struct {
	From time.Time
	To   time.Time
}

type AvailableHoursHandler cqrs.QueryHandler[AvailableHours, []Date]

type AvailableHoursReadModel interface {
	AvailableHours(ctx context.Context, from time.Time, to time.Time) ([]Date, error)
}

type availableHoursHandler struct {
	readModel AvailableHoursReadModel
}

func NewAvailableHoursHandler(
	readModel AvailableHoursReadModel,
	logger *logrus.Entry,
) AvailableHoursHandler {
	return cqrs.ApplyQueryDecorators[AvailableHours, []Date](
		availableHoursHandler{readModel: readModel},
		logger,
	)
}

func (h availableHoursHandler) Handle(ctx context.Context, query AvailableHours) (d []Date, err error) {
	if query.From.After(query.To) {
		return nil, errors.NewIncorrectInputError("date-from-after-date-to", "Date from after date to")
	}

	return h.readModel.AvailableHours(ctx, query.From, query.To)
}
