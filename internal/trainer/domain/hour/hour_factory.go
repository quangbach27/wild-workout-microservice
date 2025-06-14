package hour

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrNotFullHour = errors.New("hour should be a full hour")
	ErrPastHour    = errors.New("cannot create hour from past")
)

type TooDistantDateError struct {
	MaxWeeksInTheFutureToSet int
	ProvidedDate             time.Time
}

func (e TooDistantDateError) Error() string {
	return fmt.Sprintf(
		"schedule can be only set for next %d weeks, provided date: %s",
		e.MaxWeeksInTheFutureToSet,
		e.ProvidedDate,
	)
}

type TooEarlyHourError struct {
	MinUtcHour   int
	ProvidedTime time.Time
}

func (e TooEarlyHourError) Error() string {
	return fmt.Sprintf(
		"too early hour, min UTC hour: %d, provided time: %s",
		e.MinUtcHour,
		e.ProvidedTime,
	)
}

type TooLateHourError struct {
	MaxUtcHour   int
	ProvidedTime time.Time
}

func (e TooLateHourError) Error() string {
	return fmt.Sprintf(
		"too late hour, min UTC hour: %d, provided time: %s",
		e.MaxUtcHour,
		e.ProvidedTime,
	)
}

type HourFactory struct {
	config HourFactoryConfig
}

func NewHourFactory(config HourFactoryConfig) (HourFactory, error) {
	err := config.Validate()
	if err != nil {
		return HourFactory{}, err
	}

	return HourFactory{config: config}, nil
}

func MustNewHourFactory(config HourFactoryConfig) HourFactory {
	factory, err := NewHourFactory(config)
	if err != nil {
		panic(err)
	}

	return factory
}

func (f HourFactory) Config() HourFactoryConfig {
	return f.config
}

func (f HourFactory) NewAvailableHour(hour time.Time) (*Hour, error) {
	if err := f.validateTime(hour); err != nil {
		return nil, err
	}

	return &Hour{
		hour:         hour,
		availability: Available,
	}, nil
}

func (f HourFactory) NewNotAvailableHour(hour time.Time) (*Hour, error) {
	if err := f.validateTime(hour); err != nil {
		return nil, err
	}

	return &Hour{
		hour:         hour,
		availability: NotAvailable,
	}, nil
}

func (f HourFactory) validateTime(hour time.Time) error {
	if !hour.Round(time.Hour).Equal(hour) {
		return ErrNotFullHour
	}

	maximumTime := time.Now().AddDate(0, 0, f.config.MaxWeeksInTheFutureToSet*7)
	if hour.After(maximumTime) {
		return TooDistantDateError{
			MaxWeeksInTheFutureToSet: f.config.MaxWeeksInTheFutureToSet,
			ProvidedDate:             hour,
		}
	}

	currentHour := time.Now().Truncate(time.Hour)
	if hour.Before(currentHour) || hour.Equal(currentHour) {
		return ErrPastHour
	}

	if hour.UTC().Hour() > f.config.MaxUtcHour {
		return TooLateHourError{
			MaxUtcHour:   f.config.MaxUtcHour,
			ProvidedTime: hour,
		}
	}

	if hour.UTC().Hour() < f.config.MinUtcHour {
		return TooEarlyHourError{
			MinUtcHour:   f.config.MinUtcHour,
			ProvidedTime: hour,
		}
	}

	return nil
}
