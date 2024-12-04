package esc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Parse an import id using {project_id}:{id} structure
func parseImportID(importId string) ([]string, error) {
	result := strings.Split(importId, ":")

	if len(result) != 2 {
		return nil, fmt.Errorf(
			"Failed to parse import id. Please use the format `{project_id}:{resource_id}`",
		)
	}

	return result, nil
}

// Help to set a proper project_id and resource id for import
func resourceImport(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) ([]*schema.ResourceData, error) {
	idSlice, err := parseImportID(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("project_id", idSlice[0]); err != nil {
		return nil, err
	}
	d.SetId(idSlice[1])

	return []*schema.ResourceData{d}, nil
}
