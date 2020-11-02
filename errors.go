package solus

import (
	"encoding/json"
	"fmt"
)

type HTTPError struct {
	HttpCode int    `json:"http_code"`
	Message  string `json:"message"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.HttpCode, e.Message)
}

func newHTTPError(httpCode int, body []byte) error {
	e := HTTPError{
		HttpCode: httpCode,
	}

	if err := json.Unmarshal(body, &e); err != nil {
		e.Message = string(body)
		return e
	}

	return e
}
