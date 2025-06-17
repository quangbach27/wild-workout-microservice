package app

import "github.com/quangbach27/wild-workout-microservice/internal/trainer/app/query"

type Application struct {
	Queries Queries
}

type Queries struct {
	AvailableHours query.AvailableHoursHandler
}
