package errors

import "errors"

// These variables are used to five us access to existing
// functions in the std lib errors package. We can algo
// wrap them in custom funcionality as neede if we want,
// or mock them during testing

var (
	As = errors.As
	Is = errors.Is
)
