package adapters

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/quangbach27/wild-workout-microservice/internal/training/domain/training"
)

type TrainingRepository struct {
	db *sqlx.DB
}

type sqlGetContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func NewTrainingRepository(db *sqlx.DB) TrainingRepository {
	if db == nil {
		panic("Empty sql db")
	}

	return TrainingRepository{db: db}
}

func (repository TrainingRepository) GetTraining(
	ctx context.Context,
	trainingUUID string,
) (*training.Training, error) {
	return repository.getOrCreateTraining(ctx, repository.db, trainingUUID, false)

}

func (repository TrainingRepository) getOrCreateTraining(
	ctx context.Context,
	db sqlGetContext,
	trainingUUID string,
	forUpdate bool,
) (*training.Training, error) {
	//query := `
	//	SELECT
	//	FROM
	//`
	return nil, nil
}
