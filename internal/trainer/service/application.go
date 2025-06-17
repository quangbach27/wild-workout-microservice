package service

import (
	"context"

	"github.com/quangbach27/wild-workout-microservice/internal/trainer/adapters"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/app"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/app/query"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	db := adapters.MustNewPostgresConnection()

	hourFactoryConfig := hour.HourFactoryConfig{
		MaxWeeksInTheFutureToSet: 6,
		MinUtcHour:               12,
		MaxUtcHour:               20,
	}

	datesRepository := adapters.NewDatsRepository(db, hourFactoryConfig)
	logger := logrus.NewEntry(logrus.StandardLogger())

	return app.Application{
		Queries: app.Queries{
			AvailableHours: query.NewAvailableHoursHandler(datesRepository, logger),
		},
	}
}
