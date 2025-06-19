package adapters

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
	"go.uber.org/multierr"
	"time"
)

type HourRepository struct {
	db          *sqlx.DB
	hourFactory hour.HourFactory
}

func NewHourRepository(db *sqlx.DB, hourFactory hour.HourFactory) HourRepository {
	if db == nil {
		panic("missing db")
	}

	return HourRepository{db: db, hourFactory: hourFactory}
}

type sqlContextGetter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func (repository HourRepository) GetHour(ctx context.Context, time time.Time) (*hour.Hour, error) {
	return repository.getOrCreateHour(ctx, repository.db, time, false)
}

func (repository HourRepository) getOrCreateHour(
	ctx context.Context,
	db sqlContextGetter,
	hourTime time.Time,
	forUpdate bool,
) (*hour.Hour, error) {

	var hourModels HourModel

	query := `
		SELECT h.hour, h.availability
		FROM hours h
		WHERE h.hour = $1
	`
	if forUpdate {
		query += " FOR UPDATE"
	}

	err := db.GetContext(ctx, &hourModels, query, hourTime.UTC())
	if errors.Is(err, sql.ErrNoRows) {
		return repository.hourFactory.NewNotAvailableHour(hourTime)
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to get hour")
	}

	availability, err := hour.NewAvailabilityFromString(hourModels.Availability)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get hour from db")
	}

	domainHour, err := repository.hourFactory.UnmarshalHourFromDatabase(hourModels.Hour.Local(), availability)
	return domainHour, nil
}

const deadlockErrCode = "40P01"
const maxDeadlockRetries = 3

func (repository HourRepository) UpdateHour(
	ctx context.Context,
	hourTime time.Time,
	updateFunc func(h *hour.Hour) (*hour.Hour, error),
) (err error) {
	for attempt := 0; attempt < maxDeadlockRetries; attempt++ {
		err = repository.updateHour(ctx, hourTime, updateFunc)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == deadlockErrCode {
			// Retry on deadlock
			continue
		}

		return
	}
	return

}

func (repository HourRepository) updateHour(
	ctx context.Context,
	hourTime time.Time,
	updateFunc func(h *hour.Hour) (*hour.Hour, error),
) (err error) {
	tx, err := repository.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		err = repository.finishTransaction(err, tx)
	}()

	existingHour, err := repository.getOrCreateHour(ctx, tx, hourTime, true)
	if err != nil {
		return err
	}

	updatedHour, err := updateFunc(existingHour)
	if err != nil {
		return err
	}

	//if updatedHour.Equal(existingHour) {
	//	return nil
	//}

	if err := repository.upsertHour(tx, updatedHour); err != nil {
		return err
	}

	return err
}

// upsertHour updates hour if hour already exists in the database.
// If your doesn't exists, it's inserted.
func (repository HourRepository) upsertHour(tx *sqlx.Tx, hourToUpdate *hour.Hour) error {
	hourUtc := hourToUpdate.Hour().UTC()
	hourModel := HourModel{
		Hour:         hourUtc,
		Availability: hourToUpdate.Availability().String(),
		Date:         time.Date(hourUtc.Year(), hourUtc.Month(), hourUtc.Day(), 0, 0, 0, 0, time.UTC),
	}

	query := `
		WITH inserted_date AS (
		    INSERT INTO dates (date)
		    VALUES (:date)
		    ON CONFLICT (date) DO NOTHING 
		)
		INSERT INTO hours (date, hour, availability)
		VALUES (:date, :hour, :availability)
		ON CONFLICT (hour) DO UPDATE
		SET availability = EXCLUDED.availability
		WHERE hours.availability IS DISTINCT FROM EXCLUDED.availability
	`
	_, err := tx.NamedExec(query, hourModel)
	if err != nil {
		return errors.Wrap(err, "unable to upsert hours")
	}

	return nil
}

// finishTransaction rollbacks transaction if error is provided.
// If err is nil transaction is committed.
//
// If the rollback fails, we are using multierr library to add error about rollback failure.
// If the commit fails, commit error is returned.
func (repository HourRepository) finishTransaction(err error, tx *sqlx.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return multierr.Combine(err, rollbackErr)
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return errors.Wrap(err, "failed to commit tx")
		}

		return nil
	}
}
