package hour

var (
	Available         = Availability{value: "available"}
	NotAvailable      = Availability{value: "not_available"}
	TrainingScheduled = Availability{value: "training_scheduled"}
)

type Availability struct {
	value string
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
