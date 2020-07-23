package client

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

type problemDetails struct {
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Status   int               `json:"status"`
	Detail   string            `json:"detail"`
	Instance string            `json:"instance"`
	Fields   map[string]string `json:"fields,omitempty"`
}

func (problemDetails *problemDetails) Error() string {
	if problemDetails.Title == "" && problemDetails.Detail == "" {
		return fmt.Sprintf("Status %d", problemDetails.Status)
	}

	err := ""
	if problemDetails.Title != "" {
		err += fmt.Sprintf("%s", problemDetails.Title)
	}
	if problemDetails.Detail != "" {
		err += fmt.Sprintf(": %s", problemDetails.Detail)
	}

	if len(problemDetails.Fields) > 0 {
		keys := make([]string, 0)
		for k, _ := range problemDetails.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			err += fmt.Sprintf("\n\t%s - %s", key, problemDetails.Fields[key])
		}
	}

	return err
}

func newProblemDetailsFromReader(reader io.Reader) (*problemDetails, error) {
	var details problemDetails
	decoder := json.NewDecoder(reader)

	if err := decoder.Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}