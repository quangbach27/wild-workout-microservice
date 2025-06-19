package command

import (
	"context"
	"github.com/quangbach27/wild-workout-microservice/internal/common/cqrs"
	"github.com/quangbach27/wild-workout-microservice/internal/common/errors"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
	"github.com/sirupsen/logrus"
	"time"
)

func NewMakeHoursAvailableHandler(
	hourRepo hour.HourRepository,
	logger *logrus.Entry,
) MakeHoursAvailableHandler {
	if hourRepo == nil {
		panic("hourRepo must not be nil")
	}
	return cqrs.ApplyCommandDecorators[MakeHoursAvailableCommand](
		makeHoursAvailableHandler{hourRepo: hourRepo},
		logger,
	)
}

type MakeHoursAvailableCommand struct {
	Hours []time.Time
}

type MakeHoursAvailableHandler cqrs.CommandHandler[MakeHoursAvailableCommand]

type makeHoursAvailableHandler struct {
	hourRepo hour.HourRepository
}

func (handler makeHoursAvailableHandler) Handle(ctx context.Context, cmd MakeHoursAvailableCommand) error {
	for _, hourToUpdate := range cmd.Hours {
		if err := handler.hourRepo.UpdateHour(ctx, hourToUpdate, func(h *hour.Hour) (*hour.Hour, error) {
			if err := h.MakeAvailable(); err != nil {
				return nil, err
			}
			return h, nil
		}); err != nil {
			return errors.NewSlugError(err.Error(), "unable-to-update-availability")
		}
	}

	return nil
}
