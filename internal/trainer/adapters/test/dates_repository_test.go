package adapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/quangbach27/wild-workout-microservice/internal/trainer/adapters"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
	"github.com/stretchr/testify/require"
)

var db = adapters.MustNewPostgresConnection()
var testHourFactoryConfig = hour.HourFactoryConfig{}
var dateRepository = adapters.NewDatsRepository(db, testHourFactoryConfig)

func TestGetDate(t *testing.T) {
	fromDate, _ := time.Parse("2006-01-02", "2025-06-13")
	toDate, _ := time.Parse("2006-01-02", "2025-06-16")

	_, err := dateRepository.AvailableHours(context.Background(), fromDate, toDate)

	require.NoError(t, err)
}
