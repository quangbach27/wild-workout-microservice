package hour

import (
	"errors"
	"time"
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
	ErrNotTrainingScheduled = errors.New("no training is scheduled for this hour")
	ErrHourNotAvailable     = errors.New("hour is not available")
)

func (h *Hour) MakeAvailable() error {
	if h.availability.HasTrainingScheduled() {
		return ErrTrainingScheduled
	}

	h.availability = Available
	return nil
}

func (h *Hour) MakeUnavailable() error {
	if h.availability.HasTrainingScheduled() {
		return ErrTrainingScheduled
	}

	h.availability = Unavailable
	return nil
}

func (h *Hour) ScheduleTraining() error {
	if !h.availability.IsAvailable() {
		return ErrHourNotAvailable
	}

	if h.availability.HasTrainingScheduled() {
		return ErrTrainingScheduled
	}

	h.availability = TrainingScheduled
	return nil
}

func (h *Hour) CancelTraining() error {
	if !h.availability.HasTrainingScheduled() {
		return ErrNotTrainingScheduled
	}

	h.availability = Available
	return nil
}