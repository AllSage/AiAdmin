// Package errtypes contains custom error types
package errtypes

import (
	"fmt"
	"strings"
)

const UnknownAiAdminKeyErrMsg = "unknown AiAdmin key"
const InvalidModelNameErrMsg = "invalid model name"

// TODO: This should have a structured response from the API
type UnknownAiAdminKey struct {
	Key string
}

func (e *UnknownAiAdminKey) Error() string {
	return fmt.Sprintf("unauthorized: %s %q", UnknownAiAdminKeyErrMsg, strings.TrimSpace(e.Key))
}
