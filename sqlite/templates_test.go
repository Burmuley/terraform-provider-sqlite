package sqlite

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
	"time"
)

func Test_createTableTemplate(t *testing.T) {
	// not a real test
	// used to use it for SQL statements generation debugging
	var table struct {
		Name      string
		Created   string
		Overwrite bool
		Columns   []map[string]interface{}
	}

	table.Name = "test_table"
	table.Created = time.Now().Format(time.RFC850)
	table.Overwrite = false
	table.Columns = make([]map[string]interface{}, 0, 2)

	table.Columns = append(table.Columns, map[string]interface{}{
		"name": "tst_field1",
		"type": "INT",
		"constraints": []map[string]interface{}{
			{
				"primary_key": interface{}(true),
				"not_null":    interface{}(true),
			},
		},
	})
	table.Columns = append(table.Columns, map[string]interface{}{
		"name": "tst_field2",
		"type": "TEXT",
		"constraints": []map[string]interface{}{
			{
				"not_null": interface{}(true),
			},
		},
	})

	tmpl, err := template.New("create_table").Funcs(templateFuncMap).Parse(createTableTemplate)
	if err != nil {
		t.Error(err)
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, table)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(buf.Bytes()))
}
