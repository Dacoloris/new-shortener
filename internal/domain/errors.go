package domain

import "fmt"

type UniqueConstraintError struct {
	Err error
}

func (uc *UniqueConstraintError) Error() string {
	return fmt.Sprintf("%v", uc.Err)
}

func (uc *UniqueConstraintError) Unwrap() error {
	return uc.Err
}

func NewUniqueConstraintError(err error) error {
	return &UniqueConstraintError{
		Err: err,
	}
}
