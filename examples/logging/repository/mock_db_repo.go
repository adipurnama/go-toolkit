package repository

import (
	"context"
	"fmt"

	"github.com/adipurnama/go-toolkit/examples/logging/fakesql"
	"github.com/pkg/errors"
)

// ErrUserIDNotFound ...
type ErrUserIDNotFound int

// Error implements go-error interface.
func (e ErrUserIDNotFound) Error() string {
	return fmt.Sprintf("user id %d doesn't exist", e)
}

var (
	invalidID = 69

	// ErrInvalidUserID ...
	ErrInvalidUserID = errors.New("invalid user ID")

	userIDs = []int{
		1, 3, 4, 5, 6, 7, 8, 9, 0,
	}
)

// MockDBRepository ...
type MockDBRepository struct {
	DB *fakesql.DB
}

// FindUserByID ...
func (r *MockDBRepository) FindUserByID(ctx context.Context, id int) error {
	if id == invalidID {
		return errors.Wrapf(ErrInvalidUserID, "bad request for userID %d", id)
	}

	err := r.queryDBCall(ctx, id)

	return err
}

func (r *MockDBRepository) queryDBCall(ctx context.Context, userID int) error {
	for _, v := range userIDs {
		if v == userID {
			return nil
		}
	}

	query := fmt.Sprintf("select * from users where id = %d", userID)
	err := r.DB.QueryCtx(ctx, query)

	return errors.Wrapf(err, "r.DB.Query failed for userID %d", userID)
}
