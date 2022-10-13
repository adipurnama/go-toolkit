package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	user "github.com/adipurnama/go-toolkit/examples/echo-restapi/internal"
	"github.com/adipurnama/go-toolkit/tracer"
)

// UserRepository is dummy user repo. In reality, YOU SHOULD PREFER USING SQLC-GGENERATED CODE.
type UserRepository struct {
	db *DummyDB
}

// NewUserRepository returns user repo with dummy db.
func NewUserRepository(db *DummyDB) *UserRepository {
	return &UserRepository{db}
}

// FindUserByID finds user by specific ID inside db.
func (r *UserRepository) FindUserByID(ctx context.Context, id int) (*user.User, error) {
	ctx, span := tracer.NewSpan(ctx, tracer.SpanLvlDBQuery)
	defer span.End()

	err := r.db.QueryRowsCtx(ctx, fmt.Sprintf("select * from users where id=%d", id))
	if err != nil {
		if errors.Is(err, ErrNoRows) {
			return nil, errors.Wrap(user.ErrUserIDNotFound(id), "find user not found")
		}

		return nil, errors.Wrapf(err, "find user with id %d failed", id)
	}

	return &user.User{
		ID:    id,
		Name:  "random-user",
		Email: "user@mailaddress.com",
	}, nil
}

// CreateUser insert new user to db.
func (r *UserRepository) CreateUser(ctx context.Context, u *user.User) error {
	ctx, span := tracer.NewSpan(ctx, tracer.SpanLvlDBQuery)
	defer span.End()

	query := fmt.Sprintf("INSERT INTO users(name, email) VALUES(%s, %s) RETURNING id", u.Name, u.Email)

	id, err := r.db.QueryRowCtx(ctx, query)
	if err != nil {
		return errors.Wrapf(err, "create user for email %s failed", u.Email)
	}

	u.ID = id

	return nil
}
