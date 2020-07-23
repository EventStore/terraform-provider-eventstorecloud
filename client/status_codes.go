package client

import (
	"fmt"
	"io"
)

func translateStatusCode(status int, activity string, body io.Reader) error {
	problemDetails, err := newProblemDetailsFromReader(body)
	if err != nil {
		return err
	}

	return fmt.Errorf("error %s: %w", activity, problemDetails)
}

