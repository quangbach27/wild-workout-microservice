package training

import (
	"fmt"
	"time"
)

func (t Training) ProposedNewTime() time.Time {
	return t.proposedNewTime
}

type CantRescheduleBeforeTimeError struct {
	TrainingTime time.Time
}

func (c CantRescheduleBeforeTimeError) Error() string {
	return fmt.Sprintf(
		"can't reschedule training, not enough time before, training time: %s",
		c.TrainingTime,
	)
}
