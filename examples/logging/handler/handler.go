package handler

import (
	"context"
	"log"

	"github.com/adipurnama/go-toolkit/examples/logging/repository"
	"github.com/adipurnama/go-toolkit/examples/logging/service"
	"github.com/pkg/errors"
)

// Handler mock http handler / controller.
type Handler struct {
	S *service.Service
}

// FindUserByID ...
func (h *Handler) FindUserByID(id int) error {
	ctx := context.Background()

	err := h.S.GetUserByID(ctx, id)
	if err != nil {
		var errIDNotFound repository.ErrUserIDNotFound

		if ok := errors.As(err, &errIDNotFound); ok {
			log.Println(err)

			return err
		}

		return err
	}

	return nil
}
