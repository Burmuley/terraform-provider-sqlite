package sqlite

import (
	"bytes"
    "fmt"
    "strconv"
	"strings"
	"text/template"
)

// helper functions for templates
func templateSub(a, b int) int {
	return a - b
}

func templateMkSlice(elems ...string) []string {
	return elems
}

func templateToString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

var templateFuncMap = template.FuncMap{
	"sub":     templateSub,
	"strJoin": strings.Join,
	"mkSlice": templateMkSlice,
	"toStr":   templateToString,
	"escapeSql": escapeSQLEntity,
}

func renderTemplate(t string, data interface{}) (string, error) {
	tmpl, err := template.New("").Funcs(templateFuncMap).Parse(t)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func escapeSQLEntity(s string) string {
    ns := s
    if strings.ContainsAny(ns, `" `) {
        ns = strings.Replace(s, "\"", "\"\"", -1)
        ns = fmt.Sprintf("\"%s\"", ns)
    }
    return ns
}
