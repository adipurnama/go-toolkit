package service

import (
	"context"

	user "github.com/adipurnama/go-toolkit/example/echo-restapi/internal"
	"github.com/adipurnama/go-toolkit/example/echo-restapi/internal/repository"
	"github.com/adipurnama/go-toolkit/tracer"
)

// Service ...
type Service struct {
	repo *repository.UserRepository
}

// NewService returns new *Service instance.
func NewService(r *repository.UserRepository) *Service {
	return &Service{repo: r}
}

// FindUserByID find user by specific ID.
func (s *Service) FindUserByID(ctx context.Context, id int) (*user.User, error) {
	span := tracer.ServiceFuncSpan(ctx)
	defer span.End()

	return s.repo.FindUserByID(ctx, id)
}

// CreateUser creates new user in the system.
func (s *Service) CreateUser(ctx context.Context, u *user.User) error {
	span := tracer.ServiceFuncSpan(ctx)
	defer span.End()

	err := s.repo.CreateUser(ctx, u)
	if err != nil {
		return err
	}

	return nil
}
