package hour

import (
	"time"

	"github.com/pkg/errors"
)

type Hour struct {
	hour         time.Time
	availability Availability
}

func (h *Hour) Hour() time.Time {
	return h.hour
}

func (h *Hour) Availability() Availability {
	return h.availability
}

var (
	ErrTrainingScheduled   = errors.New("unable to modify hour, because scheduled training")
	ErrNoTrainingScheduled = errors.New("training is not scheduled")
	ErrHourNotAvailable    = errors.New("hour is not available")
)

func (h *Hour) MakeAvailable() error {
	if h.availability.HasTrainingScheduled() {
		return ErrTrainingScheduled
	}

	h.availability = Available
	return nil
}

func (h *Hour) MakeNotAvailable() error {
	if h.availability.HasTrainingScheduled() {
		return ErrTrainingScheduled
	}

	h.availability = NotAvailable
	return nil
}

func (h *Hour) ScheduleTraining() error {
	if !h.availability.IsAvailable() {
		return ErrHourNotAvailable
	}

	h.availability = TrainingScheduled
	return nil
}

func (h *Hour) CancelTraining() error {
	if !h.availability.HasTrainingScheduled() {
		return ErrNoTrainingScheduled
	}

	h.availability = Available
	return nil
}
