package command

import (
	"context"
	"github.com/quangbach27/wild-workout-microservice/internal/common/cqrs"
	"github.com/quangbach27/wild-workout-microservice/internal/common/errors"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
	"github.com/sirupsen/logrus"
	"time"
)

type MakeHoursNotAvailableCommand struct {
	Hours []time.Time
}

type MakeHoursNotAvailableHandler cqrs.CommandHandler[MakeHoursNotAvailableCommand]

func NewMakeHoursNotAvailableHandler(
	hourRepository hour.HourRepository,
	logger *logrus.Entry,
) MakeHoursNotAvailableHandler {
	if hourRepository == nil {
		panic("hourRepo must not be nil")
	}
	return cqrs.ApplyCommandDecorators[MakeHoursNotAvailableCommand](
		makeHoursNotAvailableHandler{hourRepository: hourRepository},
		logger,
	)
}

type makeHoursNotAvailableHandler struct {
	hourRepository hour.HourRepository
}

func (handler makeHoursNotAvailableHandler) Handle(ctx context.Context, cmd MakeHoursNotAvailableCommand) error {
	for _, hourToUpdate := range cmd.Hours {
		err := handler.hourRepository.UpdateHour(
			ctx,
			hourToUpdate,
			func(h *hour.Hour) (*hour.Hour, error) {
				if err := h.MakeNotAvailable(); err != nil {
					return nil, err
				}
				return h, nil
			},
		)

		if err != nil {
			return errors.NewSlugError(err.Error(), "unable-to-update-availability")
		}
	}

	return nil
}
