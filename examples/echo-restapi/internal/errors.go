package user

import "fmt"

// ErrUserIDNotFound returned when system cannot find with specific ID.
type ErrUserIDNotFound int

// Error go-error interface impl.
func (e ErrUserIDNotFound) Error() string {
	return fmt.Sprintf("user id %d doesn't exists", e)
}
