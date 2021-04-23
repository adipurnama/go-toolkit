package dto

// ResponseStatus enum for JSON response status.
type ResponseStatus string

const (
	// ResonseStatusOK means OK.
	ResonseStatusOK ResponseStatus = "OK"
)

// SuccessResponse for success response wrapper.
type SuccessResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
