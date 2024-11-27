package types

// ApiResponse represents the most basic form of response from an API.
type ApiResponse interface {
	GetStatusCode() int
	GetMessage() interface{}
	GetDetails() interface{}
}
