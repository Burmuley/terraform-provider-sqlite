package sqlite

var createTableTemplate = `
CREATE TABLE {{ .Name }} (
{{- $lenColumns := sub (len .Columns) 1 -}}
{{- range $i, $v := .Columns }}
    {{- $constr_str := "" -}}
    {{- range (index $v "constraints") }}
        {{- range $ck, $cv := . }}
            {{- if and (eq $ck "primary_key") (eq (toStr $cv) "true") }}
                {{- $constr_str = (strJoin (mkSlice "PRIMARY KEY" $constr_str) " ") }}
            {{- else if and (eq $ck "not_null") (eq (toStr $cv) "true") }}
                {{- $constr_str = (strJoin (mkSlice "NOT NULL" $constr_str) " ") }}
            {{- else if and (eq $ck "default") (gt (len (toStr $cv)) 0) }}
                {{- $constr_str = (strJoin (mkSlice "DEFAULT" $cv $constr_str) " ") }}
            {{ end -}}
        {{ end -}}
    {{ end }}
{{ index $v "name" }} {{ index $v "type" }} {{ $constr_str -}}{{- if lt $i $lenColumns -}},{{- end -}}
{{ end }}
);
`

var createIndexTemplate = `
{{- $colsStr := "" -}}
{{- if gt (len .Columns) 0 }}
{{- $cols := strJoin .Columns "," -}}
{{- $colsStr = strJoin (mkSlice " (" $cols ")") "" -}} 
{{- end -}}
CREATE {{if eq (toStr .Unique) "true"}}UNIQUE {{end}}INDEX {{ .Name }} ON {{ .Table }}{{ $colsStr }};
`
