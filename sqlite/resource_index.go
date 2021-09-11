package sqlite

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIndex() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SQLite index name.",
			},
			"table": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SQLite table name to create index on.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Index creation timestamp.",
			},
			"columns": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"unique": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
		SchemaVersion: 0,
		CreateContext: resourceIndexCreate,
		ReadContext:   resourceIndexRead,
		DeleteContext: resourceIndexDelete,
		UseJSONNumber: false,
	}
}

func resourceIndexCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	var index struct {
		Name    string
		Table   string
		Columns []string
		Unique  bool
	}
	c := m.(*sqLiteWrapper)

	index.Name = d.Get("name").(string)
	index.Table = d.Get("table").(string)
	index.Unique = d.Get("unique").(bool)
	colsRaw := d.Get("columns").([]interface{})
	index.Columns = make([]string, 0, len(colsRaw))
	for _, v := range colsRaw {
		index.Columns = append(index.Columns, templateToString(v))
	}

	query, err := renderTemplate(createIndexTemplate, index)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println(query)
	_, err = c.Exec(query)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(index.Name)
	if err := d.Set("created", time.Now().Format(time.RFC850)); err != nil {
	    return diag.FromErr(err)
    }

	return diag.Diagnostics{}
}

func resourceIndexRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	var indexName string
	type pragmaColumns struct {
		SeqNo int
		Cid   int
		Name  string
	}

	c := m.(*sqLiteWrapper)
	// SQL statements for getting table information
	// Resource Id in our case corresponds to table name
	schemaStmt := fmt.Sprintf("PRAGMA INDEX_INFO(%s);", d.Id())
	indexStmt := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='index' AND name='%s';", d.Id())
    log.Println(indexStmt)
    log.Println(schemaStmt)

	// check if table exists and get its name
	res, err := c.QueryRow(indexStmt)
	if err != nil {
		return diag.FromErr(fmt.Errorf("finding index: %w", err))
	}

	// name will be empty if table does not exist
	err = res.Scan(&indexName)
	if err != nil || len(indexName) < 1 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "index was not found in the database, re-creating...",
				Detail:   fmt.Sprintf("%s", err),
			},
		}
	}
	err = d.Set("name", indexName)
	if err != nil {
		diag.FromErr(fmt.Errorf("setting index name: %s", err))
	}

	// make up a new list of columns with the same type as defined in schema
	columns := make([]interface{}, 0)
	// query for index schema
	// result will be empty if index does not exist
	rows, err := c.Query(schemaStmt)
	if err != nil {
		return diag.FromErr(err)
	}

	for rows.Next() {
		col := pragmaColumns{}
		err = rows.Scan(&col.SeqNo, &col.Cid, &col.Name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("getting column info: %w", err))
		}
		columns = append(columns, col.Name)
	}

	// write our columns values into the resource data
	if err := d.Set("columns", columns); err != nil {
		return diag.FromErr(fmt.Errorf("reading index columns: %w", err))
	}

	return diag.Diagnostics{}
}

func resourceIndexDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sqLiteWrapper)
	query := fmt.Sprintf("DROP INDEX %s;", d.Id())
    log.Println(query)
	_, err := c.Exec(query)
	if err != nil {
		return diag.FromErr(err)
	}
	// set empty resource Id to mark it as destroyed
	// and make TF remove it from the state
	d.SetId("")
	return diag.Diagnostics{}
}
