package web_test

import (
	"net/http/httptest"
	"testing"

	"github.com/adipurnama/go-toolkit/web"
	"github.com/stretchr/testify/assert"
)

func TestGetIP(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "123.456.789.012:29302"

	// no header, default
	assert.Equal(t, "123.456.789.012", web.GetIP(r))

	// X-Real-IP
	r.Header.Set("X-Real-IP", "103.0.53.43")
	assert.Equal(t, "103.0.53.43", web.GetIP(r))

	// Forwarded
	r.Header.Set("Forwarded", "for=192.0.2.60;proto=http;by=203.0.113.43")
	assert.Equal(t, "192.0.2.60", web.GetIP(r))

	// X-Forwarded-For
	r.Header.Set("X-Forwarded-For", "127.0.0.1, 23.21.45.67")
	assert.Equal(t, "127.0.0.1", web.GetIP(r))

	// CF-Connecting-IP
	r.Header.Set("CF-Connecting-IP", "127.0.0.1, 23.21.45.67")
	assert.Equal(t, "127.0.0.1", web.GetIP(r))
}
