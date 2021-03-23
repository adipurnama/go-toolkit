package service

import (
	"context"

	"github.com/adipurnama/go-toolkit/example/logging/repository"
)

// Service mock service layer.
type Service struct {
	R *repository.MockDBRepository
}

// GetUserByID ...
func (s *Service) GetUserByID(ctx context.Context, id int) error {
	return s.getUserFromRepo(ctx, id)
}

func (s *Service) getUserFromRepo(ctx context.Context, id int) error {
	err := s.R.FindUserByID(ctx, id)
	return err
}
