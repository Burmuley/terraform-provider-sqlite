package sqlite

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SQLITE_DB_PATH", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sqlite_table": resourceTable(),
			"sqlite_index": resourceIndex(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
		ProviderMetaSchema:   map[string]*schema.Schema{},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var err error
	var diags diag.Diagnostics

	dbPath := d.Get("path").(string)
	if len(dbPath) < 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "parameter 'path' can not be empty",
		})
		return nil, diags
	}

	sqlW := NewSqLiteWrapper()
	err = sqlW.Open(dbPath)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("error opening the database '%s'", dbPath),
			Detail:   fmt.Sprint(err),
		})
		return nil, diags
	}

	return sqlW, diags
}
