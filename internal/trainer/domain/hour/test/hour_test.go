package hour_test

import (
	"testing"

	"github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testHourFactory = hour.MustNewHourFactory(hour.HourFactoryConfig{
	MaxWeeksInTheFutureToSet: 100,
	MinUtcHour:               0,
	MaxUtcHour:               24,
})

func TestHour_MakeNotAvailable(t *testing.T) {
	h, err := testHourFactory.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	require.NoError(t, h.MakeNotAvailable())
	assert.False(t, h.Availability().IsAvailable())
}

func TestHour_MakeNotAvailable_with_scheduled_training(t *testing.T) {
	h := newHourWithScheduledTraining(t)

	err := h.MakeNotAvailable()
	assert.ErrorIs(t, err, hour.ErrTrainingScheduled)
}

func TestHour_MakeAvailable(t *testing.T) {
	h, err := testHourFactory.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	require.NoError(t, h.MakeNotAvailable())

	require.NoError(t, h.MakeAvailable())
	assert.True(t, h.Availability().IsAvailable())
}

func TestHour_MakeAvailable_with_scheduled_training(t *testing.T) {
	t.Parallel()
	h := newHourWithScheduledTraining(t)

	assert.Equal(t, hour.ErrTrainingScheduled, h.MakeAvailable())
}

func TestHour_ScheduleTraining(t *testing.T) {
	t.Parallel()
	h, err := testHourFactory.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	require.NoError(t, h.ScheduleTraining())

	assert.True(t, h.Availability().HasTrainingScheduled())
	assert.False(t, h.Availability().IsAvailable())
}

func TestHour_ScheduleTraining_with_not_available(t *testing.T) {
	t.Parallel()
	h := newNotAvailableHour(t)
	assert.Equal(t, hour.ErrHourNotAvailable, h.ScheduleTraining())
}

func TestHour_CancelTraining(t *testing.T) {
	t.Parallel()
	h := newHourWithScheduledTraining(t)

	require.NoError(t, h.CancelTraining())

	assert.False(t, h.Availability().HasTrainingScheduled())
	assert.True(t, h.Availability().IsAvailable())
}

func TestHour_CancelTraining_no_training_scheduled(t *testing.T) {
	t.Parallel()
	h, err := testHourFactory.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	err = h.CancelTraining()
	assert.ErrorIs(t, err, hour.ErrNoTrainingScheduled)
}

func newHourWithScheduledTraining(t *testing.T) *hour.Hour {
	h, err := testHourFactory.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)
	assert.Equal(t, hour.Available, h.Availability())

	err = h.ScheduleTraining()
	require.NoError(t, err)
	assert.Equal(t, hour.TrainingScheduled, h.Availability())

	return h
}

func newNotAvailableHour(t *testing.T) *hour.Hour {
	h, err := testHourFactory.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	err = h.MakeNotAvailable()
	require.NoError(t, err)

	return h
}
