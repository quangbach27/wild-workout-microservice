package hour_test

import (
	"testing"

	"github.com/quangbach27/wild-workout-microservice/internal/trainer/domain/hour"
	"github.com/stretchr/testify/assert"
)

func TestFactoryConfig_Validate(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		Config      hour.HourFactoryConfig
		ExpectedErr string
	}{
		{
			Name: "valid",
			Config: hour.HourFactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               10,
				MaxUtcHour:               12,
			},
			ExpectedErr: "",
		},
		{
			Name: "equal_min_and_max_hour",
			Config: hour.HourFactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               12,
				MaxUtcHour:               12,
			},
			ExpectedErr: "",
		},
		{
			Name: "min_hour_after_max_hour",
			Config: hour.HourFactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               13,
				MaxUtcHour:               12,
			},
			ExpectedErr: "MinUtcHour (13) can't be after MaxUtcHour (12)",
		},
		{
			Name: "zero_max_weeks",
			Config: hour.HourFactoryConfig{
				MaxWeeksInTheFutureToSet: 0,
				MinUtcHour:               10,
				MaxUtcHour:               12,
			},
			ExpectedErr: "MaxWeeksInTheFutureToSet should be greater than 1, but is 0",
		},
		{
			Name: "sub_zero_min_hour",
			Config: hour.HourFactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               -1,
				MaxUtcHour:               12,
			},
			ExpectedErr: "MinUtcHour should be value between 0 and 24, but is -1",
		},
		{
			Name: "sub_zero_max_hour",
			Config: hour.HourFactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               10,
				MaxUtcHour:               -1,
			},
			ExpectedErr: "MinUtcHour should be value between 0 and 24, but is -1; MinUtcHour (10) can't be after MaxUtcHour (-1)",
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			err := c.Config.Validate()

			if c.ExpectedErr != "" {
				assert.EqualError(t, err, c.ExpectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
