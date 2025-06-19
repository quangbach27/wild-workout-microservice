package app

import (
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/app/command"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/app/query"
)

type Application struct {
	Queries  Queries
	Commands Commands
}

type Queries struct {
	AvailableHours query.AvailableHoursHandler
}
type Commands struct {
	MakeHoursAvailable command.MakeHoursAvailableHandler
}
