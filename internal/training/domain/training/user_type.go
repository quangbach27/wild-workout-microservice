package training

import (
	"fmt"
	commonErrors "github.com/quangbach27/wild-workout-microservice/internal/common/errors"
)

type UserType struct {
	s string
}

func (u UserType) IsZero() bool {
	return u == UserType{}
}

func (u UserType) String() string {
	return u.s
}

var (
	Trainer  = UserType{"trainer"}
	Attendee = UserType{"attendee"}
)

func NewUserTypeFromString(userType string) (UserType, error) {
	switch userType {
	case "trainer":
		return Trainer, nil
	case "attendee":
		return Attendee, nil
	}

	return UserType{}, commonErrors.NewIncorrectInputError(
		fmt.Sprintf("invalid '%s' role", userType),
		"invalid-role",
	)
}
