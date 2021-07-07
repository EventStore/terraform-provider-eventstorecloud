package esc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var defaultTokenStore = filepath.Join(os.Getenv("HOME"), ".esctf", "tokens")

func init() {
	schema.DescriptionKind = schema.StringMarkdown

	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Deprecated != "" {
			desc += " " + s.Deprecated
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"url": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("ESC_URL", "https://api.eventstore.cloud"),
				},

				"token": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("ESC_TOKEN", ""),
					Sensitive:   true,
				},

				"organization_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("ESC_ORG_ID", ""),
				},

				"token_store": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("ESC_TOKEN_STORE", defaultTokenStore),
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"eventstorecloud_project": dataSourceProject(),
			},

			ResourcesMap: map[string]*schema.Resource{
				"eventstorecloud_project":          resourceProject(),
				"eventstorecloud_network":          resourceNetwork(),
				"eventstorecloud_peering":          resourcePeering(),
				"eventstorecloud_managed_cluster":  resourceManagedCluster(),
				"eventstorecloud_scheduled_backup": resourceScheduledBackup(),
				"eventstorecloud_integration":      resourceIntegration(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}
func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := &client.Config{
			URL:          d.Get("url").(string),
			RefreshToken: d.Get("token").(string),
			TokenStore:   d.Get("token_store").(string),
		}

		c, err := client.New(config)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return &providerContext{
			organizationId: d.Get("organization_id").(string),
			client:         c,
		}, nil
	}
}

type providerContext struct {
	organizationId string
	client         *client.Client
}

// Networks, Peerings
var fastResourceTimeout = &schema.ResourceTimeout{
	Create: schema.DefaultTimeout(3 * time.Minute),
	Read:   schema.DefaultTimeout(30 * time.Second),
	Update: schema.DefaultTimeout(3 * time.Minute),
	Delete: schema.DefaultTimeout(3 * time.Minute),
}

// Clusters
var slowResourceTimeout = &schema.ResourceTimeout{
	Create: schema.DefaultTimeout(10 * time.Minute),
	Read:   schema.DefaultTimeout(30 * time.Second),
	Update: schema.DefaultTimeout(10 * time.Minute),
	Delete: schema.DefaultTimeout(10 * time.Minute),
}

var validProviders = []string{"aws", "gcp", "azure"}
var validServerVersions = []string{"20.6", "20.10"}
var validTopologies = []string{"single-node", "three-node-multi-zone"}
var validInstanceTypes = []string{"F1", "C4", "M8", "M16", "M32", "M64", "M128"}
var validDiskTypes = []string{"GP2", "SSD", "PREMIUM-SSD-LRS"}
var validProjectionLevels = []string{"off", "system", "user"}
