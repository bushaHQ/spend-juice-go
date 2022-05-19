package juice

import (
	"fmt"
	"strings"
)

func (e Error) Error() string {
	errorBuilder := strings.Builder{}
	errorBuilder.WriteString(e.Message + " ")
	if e.Errors != nil {
		errorBuilder.WriteString(fmt.Sprintf("%v", e.Errors) + "; ")
	}
	return strings.ToLower(strings.Trim(errorBuilder.String(), ";. "))
}

type Error struct {
	Errors  interface{} `json:"errors"`
	Message string      `json:"message"`
}
