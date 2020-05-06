package web

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHeaderFromContext(t *testing.T) {
	ctx := context.Background()
	header := http.Header{}
	header.Set("reqid", "ID")

	t.Run("header found", func(t *testing.T) {
		ctx = context.WithValue(ctx, ContextKeyHeader, header)

		got := HeaderFromContext(ctx)
		if got["Reqid"][0] != header["Reqid"][0] {
			t.Errorf("should return previously set header")
		}
	})
}

func TestContextKey(t *testing.T) {
	key := ContextKey("Hello")
	assert.Equal(t, "web context: Hello", key.String())
}
