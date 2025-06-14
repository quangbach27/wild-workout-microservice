package hour_test

import (
	"testing"
	"time"

	"github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAvailableHour(t *testing.T) {
	t.Parallel()
	h, err := testHourFactory.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	assert.True(t, h.Availability().IsAvailable())
}

func TestNewAvailableHour_not_full_hour(t *testing.T) {
	t.Parallel()
	constructorTime := trainingHourWithMinutes(13)

	_, err := testHourFactory.NewAvailableHour(constructorTime)
	assert.Equal(t, hour.ErrNotFullHour, err)
}

func TestNewAvailableHour_too_distant_date(t *testing.T) {
	t.Parallel()
	maxWeeksInFuture := 1

	factory := hour.MustNewHourFactory(hour.HourFactoryConfig{
		MaxWeeksInTheFutureToSet: maxWeeksInFuture,
		MinUtcHour:               0,
		MaxUtcHour:               0,
	})

	constructorTime := time.Now().Truncate(time.Hour*24).AddDate(0, 0, maxWeeksInFuture*7+1)

	_, err := factory.NewAvailableHour(constructorTime)
	assert.Equal(
		t,
		hour.TooDistantDateError{
			MaxWeeksInTheFutureToSet: maxWeeksInFuture,
			ProvidedDate:             constructorTime,
		},
		err,
	)
}

func TestNewAvailableHour_past_date(t *testing.T) {
	t.Parallel()
	pastHour := time.Now().Truncate(time.Hour).Add(-time.Hour)
	_, err := testHourFactory.NewAvailableHour(pastHour)
	assert.Equal(t, hour.ErrPastHour, err)

	currentHour := time.Now().Truncate(time.Hour)
	_, err = testHourFactory.NewAvailableHour(currentHour)
	assert.Equal(t, hour.ErrPastHour, err)
}

func TestNewAvailableHour_too_early_hour(t *testing.T) {
	t.Parallel()
	factory := hour.MustNewHourFactory(hour.HourFactoryConfig{
		MaxWeeksInTheFutureToSet: 10,
		MinUtcHour:               12,
		MaxUtcHour:               18,
	})

	// we are using next day, to be sure that provided hour is not in the past
	currentTime := time.Now().AddDate(0, 0, 1)

	tooEarlyHour := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		factory.Config().MinUtcHour-1, 0, 0, 0,
		time.UTC,
	)

	_, err := factory.NewAvailableHour(tooEarlyHour)
	assert.Equal(
		t,
		hour.TooEarlyHourError{
			MinUtcHour:   factory.Config().MinUtcHour,
			ProvidedTime: tooEarlyHour,
		},
		err,
	)
}

func TestNewAvailableHour_too_late_hour(t *testing.T) {
	t.Parallel()
	factory := hour.MustNewHourFactory(hour.HourFactoryConfig{
		MaxWeeksInTheFutureToSet: 10,
		MinUtcHour:               12,
		MaxUtcHour:               18,
	})

	// we are using next day, to be sure that provided hour is not in the past
	currentTime := time.Now().AddDate(0, 0, 1)

	tooEarlyHour := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		factory.Config().MaxUtcHour+1, 0, 0, 0,
		time.UTC,
	)

	_, err := factory.NewAvailableHour(tooEarlyHour)
	assert.Equal(
		t,
		hour.TooLateHourError{
			MaxUtcHour:   factory.Config().MaxUtcHour,
			ProvidedTime: tooEarlyHour,
		},
		err,
	)
}

func TestHour_Time(t *testing.T) {
	t.Parallel()
	expectedTime := validTrainingHour()

	h, err := testHourFactory.NewAvailableHour(expectedTime)
	require.NoError(t, err)

	assert.Equal(t, expectedTime, h.Hour())
}

func validTrainingHour() time.Time {
	tomorrow := time.Now().Add(time.Hour * 24)

	return time.Date(
		tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		testHourFactory.Config().MinUtcHour, 0, 0, 0,
		time.UTC,
	)
}

func trainingHourWithMinutes(minute int) time.Time {
	tomorrow := time.Now().Add(time.Hour * 24)

	return time.Date(
		tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		testHourFactory.Config().MaxUtcHour, minute, 0, 0,
		time.UTC,
	)
}
