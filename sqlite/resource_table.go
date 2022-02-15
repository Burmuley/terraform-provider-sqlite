package sqlite

import (
	"context"
	"fmt"
    "log"
    "reflect"
    "time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

//func checkColumnChanges(old, new map[string]interface{}) bool {
//    // check if type has changed
//    if old["type"].(string) != new["type"].(string) {
//        return false
//    }
//
//    oldConstr, ok := old["constraints"]
//
//}

func resourceTable() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "SQLite table name.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Table creation timestamp.",
			},
            "last_updated": {
                Type:        schema.TypeString,
                Computed:    true,
                Optional:    true,
                Description: "Table lastr updated timestamp.",
            },
			"column": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
							Description: "Column name.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
							Description: "Column data type.",
							ValidateFunc: validation.StringInSlice([]string{
								"INTEGER", "TEXT", "BLOB", "REAL", "NUMERIC",
							}, false),
						},
						"constraints": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							ForceNew:    false,
							Description: "The list of column constraints.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"primary_key": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: false,
										Default:  false,
									},
									"not_null": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: false,
										Default:  false,
									},
									"default": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: false,
										Default:  nil,
									},
								},
							},
						},
					},
				},
			},
		},
		SchemaVersion: 0,
		CreateContext: resourceTableCreate,
		ReadContext:   resourceTableRead,
		DeleteContext: resourceTableDelete,
		UpdateContext: resourceTableUpdate,
		CustomizeDiff: customdiff.All(
		    customdiff.ForceNewIfChange("column", func(ctx context.Context, old, new, meta interface{}) bool {
                // always check all old columns, if Type or Constraints changed -> return true
		        // check if column was deleted (len(new)<len(old))
		        oldColumns := old.([]interface{})
		        newColumns := new.([]interface{})

		        // check if old columns didn't change sensitive fields
		        // like `type` and `constraints`
		        if len(newColumns) >= len(oldColumns) {
                    for i := range oldColumns {
                        if !reflect.DeepEqual(oldColumns[i], newColumns[i]) {
                            return true
                        }
                    }
                }

                return false
            }),
        ),
		Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },
		UseJSONNumber: false,
	}
}

func resourceTableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	var table struct {
		Name    string
		Created string
		Columns []map[string]interface{}
	}

	c := m.(*sqLiteWrapper)
	table.Name = d.Get("name").(string)
	table.Created = d.Get("created").(string)
	columns := d.Get("column").([]interface{})
	table.Columns = make([]map[string]interface{}, 0, len(columns))

	for _, column := range columns {
		colMap := column.(map[string]interface{})
		var constrMaps []map[string]interface{}
		if constraints, ok := colMap["constraints"]; ok {
			for _, c := range constraints.([]interface{}) {
				constrMaps = append(constrMaps, c.(map[string]interface{}))
			}
			colMap["constraints"] = constrMaps
		}
		table.Columns = append(table.Columns, colMap)
	}

	query, err := renderTemplate(createTableTemplate, table)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println(query)

	_, err = c.Exec(query)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(table.Name)
	if err := d.Set("created", time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(fmt.Errorf("set created: %w", err))
	}

	return resourceTableRead(ctx, d, m)
}

func resourceTableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// list of columns returned by `PRAGMA TABLE_INFO` query
	type pragmaColumns struct {
		Cid          int
		Name         string
		DataType     string
		NotNull      int
		DefaultValue interface{}
		PrimaryKey   int
	}

	var err error
	var tableName string

	c := m.(*sqLiteWrapper)
	// SQL statements for getting table information
	// Resource Id in our case correspons to table name
    SchemaStmt := fmt.Sprintf("PRAGMA TABLE_INFO(%s);", escapeSQLEntity(d.Id()))
    TableStmt := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", d.Id())
    log.Println(TableStmt)
    log.Println(SchemaStmt)

	// check if table exists and get its name
	res, err := c.QueryRow(TableStmt)
	if err != nil {
		return diag.FromErr(err)
	}

	// name will be empty if table does not exist
	err = res.Scan(&tableName)
	if err != nil || len(tableName) < 1 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "table was not found in the database, re-creating...",
				Detail:   fmt.Sprintf("%s", err),
			},
		}
	}
	err = d.Set("name", tableName)
	if err != nil {
		diag.FromErr(fmt.Errorf("error setting table name: %s", err))
	}

	// make up a new list of columns with the same type as defined in schema
	columns := make([]map[string]interface{}, 0)
	// query for table schema
	// result will be empty if table does not exist
	rows, err := c.Query(SchemaStmt)
	if err != nil {
		return diag.FromErr(err)
	}
	// iterate over result and retrieve columns configuration
	// see details: https://www.sqlite.org/pragma.html#pragma_table_info
	for rows.Next() {
		col := pragmaColumns{}
		constraints := make([]map[string]interface{}, 0)
		err = rows.Scan(&col.Cid, &col.Name, &col.DataType, &col.NotNull, &col.DefaultValue, &col.PrimaryKey)
		if err != nil {
			return diag.FromErr(err)
		}
		cConstr := make(map[string]interface{}, 0)
		if col.PrimaryKey > 0 {
			cConstr["primary_key"] = true
		}
		if col.NotNull > 0 {
			cConstr["not_null"] = true
		}
		if col.DefaultValue != nil {
			cConstr["default"] = col.DefaultValue
		}
		if len(cConstr) > 0 {
			constraints = append(constraints, cConstr)
		}

		dCol := map[string]interface{}{
			"name": col.Name,
			"type": col.DataType,
		}
		dCol["constraints"] = constraints

		columns = append(columns, dCol)
	}

	// write our columns values into the resource data
	if err := d.Set("column", columns); err != nil {
		return diag.FromErr(fmt.Errorf("reading table columns: %w", err))
	}

	return diag.Diagnostics{}
}

func resourceTableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sqLiteWrapper)
	query := fmt.Sprintf("DROP TABLE %s;", escapeSQLEntity(d.Id()))
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

func resourceTableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    c := m.(*sqLiteWrapper)
    queries := make([]string, 0)

    if d.HasChange("name") {
        oldName, newName := d.GetChange("name")
        tableRenameQuery := fmt.Sprintf("ALTER TABLE '%s' RENAME TO '%s';", oldName, newName)
        log.Println(tableRenameQuery)
        queries = append(queries, tableRenameQuery)
    }

    //if d.HasChange("column") {
    //    // something to do with columns
    //}


    for _, query := range queries {
        _, err := c.Exec(query)
        if err != nil {
            return diag.FromErr(err)
        }
    }

    if len(queries) > 0 {
        d.Set("last_updated", time.Now().Format(time.RFC850))
        d.SetId(d.Get("name").(string))
    }

    return resourceTableRead(ctx, d, m)
}
