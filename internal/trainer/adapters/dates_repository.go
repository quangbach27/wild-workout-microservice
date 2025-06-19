package adapters

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/app/query"
	domainHour "github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
)

type RawDateModels struct {
	Date  time.Time       `db:"date"`
	Hours json.RawMessage `db:"hours"`
}

type DatesRepository struct {
	db                *sqlx.DB
	hourFactoryConfig domainHour.HourFactoryConfig
}

func NewDatsRepository(db *sqlx.DB, hourFactoryConfig domainHour.HourFactoryConfig) DatesRepository {
	if db == nil {
		panic("missing db")
	}

	return DatesRepository{db: db, hourFactoryConfig: hourFactoryConfig}
}

func (repository DatesRepository) AvailableHours(ctx context.Context, from time.Time, to time.Time) ([]query.Date, error) {
	queryStr := `
		SELECT 
			d.date,
			json_agg(
				json_build_object(
					'hour', to_char(h.hour, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
					'availability', h.availability
				)
			) AS hours
		FROM dates d
			LEFT JOIN hours h 
				ON h.date = d.date
		WHERE d.date BETWEEN $1 AND $2
		GROUP BY d.date
		ORDER BY d.date
	`

	var rawDateModels []RawDateModels

	err := repository.db.SelectContext(ctx, &rawDateModels, queryStr, from, to)
	if err != nil {
		return nil, err
	}

	// Convert raw date models to structured dateModels
	dateModels, err := unmarshalDateModels(rawDateModels)
	if err != nil {
		return nil, err
	}

	queryDates, err := repository.dateModelsToApp(dateModels, from, to)
	if err != nil {
		return nil, err
	}
	return queryDates, nil
}

func (repository DatesRepository) dateModelsToApp(dateModels []DateModel, from, to time.Time) ([]query.Date, error) {
	dayMaps := make(map[time.Time]DateModel)
	for _, dateModel := range dateModels {
		dayMaps[dateModel.Date] = dateModel
	}

	var queryDates []query.Date
	for day := from; !day.After(to); day = day.AddDate(0, 0, 1) {
		var queryDate query.Date
		var err error

		if dateModel, ok := dayMaps[day]; ok {
			queryDate, err = repository.convertAndSetDefaultAvailability(dateModel)
		} else {
			queryDate, err = repository.convertAndSetDefaultAvailability(DateModel{
				Date:  day,
				Hours: nil,
			})
		}

		if err != nil {
			return nil, err
		}
		queryDates = append(queryDates, queryDate)
	}

	return queryDates, nil
}

func (repository DatesRepository) convertAndSetDefaultAvailability(dateModel DateModel) (query.Date, error) {
	day := dateModel.Date
	date := query.Date{
		Date:  dateModel.Date,
		Hours: []query.Hour{},
	}

	hourMap := make(map[time.Time]HourModel)
	for _, h := range dateModel.Hours {
		hourMap[h.Hour] = h
	}

	for h := repository.hourFactoryConfig.MinUtcHour; h <= repository.hourFactoryConfig.MaxUtcHour; h++ {
		hour := time.Date(day.Year(), day.Month(), day.Day(), h, 0, 0, 0, time.UTC)

		if hourModel, ok := hourMap[hour]; ok {
			availability, err := domainHour.NewAvailabilityFromString(hourModel.Availability)
			if err != nil {
				return query.Date{}, err
			}

			if availability.IsAvailable() && !date.HasFreeHours {
				date.HasFreeHours = true
			}
			newHour := query.Hour{
				Hour:                 hourModel.Hour,
				HasTrainingScheduled: availability.HasTrainingScheduled(),
				Available:            availability.IsAvailable(),
			}
			date.Hours = append(date.Hours, newHour)
			continue
		}

		newHour := query.Hour{
			Available: false,
			Hour:      hour,
		}

		date.Hours = append(date.Hours, newHour)
	}

	return date, nil
}

func unmarshalDateModels(rawModels []RawDateModels) ([]DateModel, error) {
	var dateModels []DateModel

	for _, raw := range rawModels {
		var hours []HourModel
		if len(raw.Hours) > 0 {
			if err := json.Unmarshal(raw.Hours, &hours); err != nil {
				return nil, err
			}
		}

		dateModels = append(dateModels, DateModel{
			Date:  raw.Date,
			Hours: hours,
		})
	}

	return dateModels, nil
}
