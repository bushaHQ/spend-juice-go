package juice

import (
	"strings"
)

func (e Error) Error() string {
	errorBuilder := strings.Builder{}
	errorBuilder.WriteString(e.Message + " ")
	if e.Errors.Amount != nil {
		errorBuilder.WriteString("(amount" + e.Errors.Amount[0] + ")" + "; ")
	}
	if e.Errors.Domain != nil {
		errorBuilder.WriteString(e.Errors.Domain[0] + "; ")
	}
	if e.Errors.Message != "" {
		errorBuilder.WriteString(e.Errors.Message + "; ")
	}
	if e.Errors.PhoneNumber != nil {
		errorBuilder.WriteString(e.Errors.PhoneNumber[0] + "; ")
	}
	if e.Errors.JuiceUserId != nil {
		errorBuilder.WriteString("(" + e.Errors.JuiceUserId[0] + " juice user id)" + " ")
	}
	return strings.ToLower(strings.Trim(errorBuilder.String(), ";. "))
}

type Error struct {
	Errors struct {
		Message     string   `json:"message"`
		Amount      []string `json:"amount"`
		Domain      []string `json:"domain"`
		PhoneNumber []string `json:"phone_number"`
		JuiceUserId []string `json:"juice_user_id"`
	} `json:"errors"`
	Message string `json:"message"`
}
