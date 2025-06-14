package hour_test

import (
	"testing"

	"github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
	"github.com/stretchr/testify/assert"
)

func TestHourAvailable(t *testing.T) {
	available := hour.Available

	t.Run("should be available", func(t *testing.T) {
		assert.True(t, available.IsAvailable())
	})

	t.Run("should not have training scheduled", func(t *testing.T) {
		assert.False(t, available.HasTrainingScheduled())
	})

	t.Run("should not be zero", func(t *testing.T) {
		assert.False(t, available.IsZero())
	})
}

func TestHourNotAvailable(t *testing.T) {
	notAvailable := hour.NotAvailable

	t.Run("should not be available", func(t *testing.T) {
		assert.False(t, notAvailable.IsAvailable())
	})

	t.Run("should not have training scheduled", func(t *testing.T) {
		assert.False(t, notAvailable.HasTrainingScheduled())
	})

	t.Run("should not be zero", func(t *testing.T) {
		assert.False(t, notAvailable.IsZero())
	})
}
func TestHourHasTrainingScheduled(t *testing.T) {
	trainingScheduled := hour.TrainingScheduled

	t.Run("should not be available", func(t *testing.T) {
		assert.False(t, trainingScheduled.IsAvailable())
	})

	t.Run("should have training scheduled", func(t *testing.T) {
		assert.True(t, trainingScheduled.HasTrainingScheduled())
	})

	t.Run("should not be zero", func(t *testing.T) {
		assert.False(t, trainingScheduled.IsZero())
	})
}
