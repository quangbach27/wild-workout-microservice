package training

import (
	"github.com/pkg/errors"
	commonerrors "github.com/quangbach27/wild-workout-microservice/internal/common/errors"
	"time"
)

type Training struct {
	uuid string

	time  time.Time
	notes string

	proposedNewTime time.Time

	canceled bool
}

func NewTraining(uuid string, userUUID string, trainingTime time.Time) (*Training, error) {
	if uuid == "" {
		return nil, errors.New("empty training uuid")
	}
	if userUUID == "" {
		return nil, errors.New("empty userUUID")
	}
	if trainingTime.IsZero() {
		return nil, errors.New("zero training time")
	}

	return &Training{
		uuid: uuid,
		time: trainingTime,
	}, nil
}

func UnmarshalTrainingFromDatabase(
	uuid string,
	userUUID string,
	trainingTime time.Time,
	notes string,
	canceled bool,
	proposedNewTime time.Time,
) (*Training, error) {
	tr, err := NewTraining(uuid, userUUID, trainingTime)
	if err != nil {
		return nil, err
	}

	tr.notes = notes
	tr.proposedNewTime = proposedNewTime
	tr.canceled = canceled

	return tr, nil
}

func (t Training) UUID() string {
	return t.uuid
}

func (t Training) Time() time.Time {
	return t.time
}

var ErrNoteTooLong = commonerrors.NewIncorrectInputError("Note too long", "note-too-long")

func (t *Training) UpdateNotes(notes string) error {
	if len(notes) > 1000 {
		return errors.WithStack(ErrNoteTooLong)
	}

	t.notes = notes
	return nil
}

func (t Training) Notes() string {
	return t.notes
}
