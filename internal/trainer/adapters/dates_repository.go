package adapters

import (
	"context"
	"encoding/json"
	"fmt"
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

	dateModels, err := unmarshalDateModels(rawDateModels)
	if err != nil {
		return nil, err
	}

	var queryDates []query.Date

	for _, dateModel := range dateModels {
		currentDate := dateModel.Date
		if currentDate.After(to) {
			break
		} else if currentDate.Equal(from) || dateModel.Date.Equal(to) {
			// convert dateModel to App
			// if len()
			queryDate := repository.convertAndSetDefaultAvailability(dateModel)
			queryDates = append(queryDates, queryDate)
		} else {
			// loop from From_Date to currentDate
			for day := from; day.Before(currentDate); day = day.AddDate(0, 0, 1) {
				queryDate := repository.convertAndSetDefaultAvailability(DateModel{
					Date:  day,
					Hours: nil,
				})
				queryDates = append(queryDates, queryDate)
			}
		}
		from = currentDate.AddDate(0, 0, 1)
	}
	fmt.Println("queryDates: ", queryDates)

	return queryDates, err
}

func (repository DatesRepository) convertAndSetDefaultAvailability(dateModel DateModel) query.Date {
	rangeHour := repository.hourFactoryConfig.MaxUtcHour - repository.hourFactoryConfig.MinUtcHour
	date := query.Date{
		Date:  dateModel.Date,
		Hours: make([]query.Hour, rangeHour),
	}

	hourMap := make(map[time.Time]HourModel)
	for _, h := range dateModel.Hours {
		hourMap[h.Hour] = h
	}

	for h := repository.hourFactoryConfig.MinUtcHour; h <= repository.hourFactoryConfig.MaxUtcHour; h++ {
		hour := time.Date(dateModel.Date.Year(), dateModel.Date.Month(), dateModel.Date.Day(), h, 0, 0, 0, time.UTC)

		if hourModel, ok := hourMap[hour]; ok {
			availability, err := domainHour.NewAvailabilityFromString(hourModel.Availability)
			if err != nil {
				continue
			}

			if availability.IsAvailable() && !date.HasFreeHours {
				date.HasFreeHours = true
			}
			continue
		}

		newHour := query.Hour{
			Available: false,
			Hour:      hour,
		}

		date.Hours = append(date.Hours, newHour)
	}

	return date
}

func unmarshalDateModels(rawModels []RawDateModels) ([]DateModel, error) {
	dateModels := make([]DateModel, len(rawModels))

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
