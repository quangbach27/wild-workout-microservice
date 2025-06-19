package adapters

import "time"

type DateModel struct {
	Date  time.Time   `db:"date"`
	Hours []HourModel // Mannually map it
}

type HourModel struct {
	Hour         time.Time `db:"hour" json:"hour"`
	Availability string    `db:"availability" json:"availability"`
	Date         time.Time `db:"date"`
}
