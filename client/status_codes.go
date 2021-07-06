package client

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"io"
)

func translateStatusCode(status int, activity string, body io.Reader) diag.Diagnostics {
	problemDetails, err := newProblemDetailsFromReader(body)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Errorf("error %s: %w", activity, problemDetails)
}
