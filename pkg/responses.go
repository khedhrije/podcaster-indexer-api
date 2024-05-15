package pkg

// ErrorJSON represents the structure for error messages in JSON responses.
type ErrorJSON struct {
	Error string `json:"error" description:"error message"`
}
