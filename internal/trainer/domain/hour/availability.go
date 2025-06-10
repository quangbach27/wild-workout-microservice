package hour

import "github.com/pkg/errors"

var (
	Available         = Availability{value: "available"}
	Unavailable       = Availability{value: "unavailable"}
	TrainingScheduled = Availability{value: "training_scheduled"}
)

// Availability is enum
type Availability struct {
	value string
}

func NewAvailabilityFromString(value string) (Availability, error) {
	switch value {
	case Available.value:
		return Available, nil
	case Unavailable.value:
		return Unavailable, nil
	case TrainingScheduled.value:
		return TrainingScheduled, nil
	default:
		return Availability{}, errors.Errorf("invalid availability value: %s", value)
	}
}

func (a Availability) String() string {
	return a.value
}

func (a Availability) IsAvailable() bool {
	return a.value == Available.value
}

func (a Availability) HasTrainingScheduled() bool {
	return a.value == TrainingScheduled.value
}

func (a Availability) IsZero() bool {
	return a == Availability{}
}
