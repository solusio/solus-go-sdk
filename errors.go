package solus

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	HttpCode int    `json:"http_code"`
	Message  string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.HttpCode, e.Message)
}

func wrapError(httpCode int, body []byte) error {
	e := Error{
		HttpCode: httpCode,
	}

	if err := json.Unmarshal(body, &e); err != nil {
		e.Message = string(body)
		return e
	}

	return e
}
