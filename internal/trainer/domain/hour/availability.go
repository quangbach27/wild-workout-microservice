package hour

import "github.com/pkg/errors"

var (
	Available         = Availability{value: "available"}
	NotAvailable      = Availability{value: "not_available"}
	TrainingScheduled = Availability{value: "training_scheduled"}
)

var availabilityValues = []Availability{
	Available,
	NotAvailable,
	TrainingScheduled,
}

type Availability struct {
	value string
}

func NewAvailabilityFromString(availabilityStr string) (Availability, error) {
	for _, availability := range availabilityValues {
		if availability.String() == availabilityStr {
			return availability, nil
		}
	}
	return Availability{}, errors.Errorf("unknown '%s' availability", availabilityStr)
}

func (a Availability) IsZero() bool {
	return a == Availability{}
}

func (a Availability) String() string {
	return a.value
}

func (a Availability) IsAvailable() bool {
	return a == Available
}

func (a Availability) HasTrainingScheduled() bool {
	return a == TrainingScheduled
}
