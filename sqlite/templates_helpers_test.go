package sqlite

import "testing"

func Test_escapeTableName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
        {
            name: "Name with space",
            args: args{s: "table name"},
            want: `"table name"`,
        },
        {
            name: "Name with double quotes",
            args: args{s: `table"name`},
            want: `"table""name"`,
        },
        {
            name: "Name with double quotes and space",
            args: args{s: `table"name two`},
            want: `"table""name two"`,
        },
        {
            name: "Name with no quotes and spaces",
            args: args{s: "table_name"},
            want: "table_name",
        },
	}
    for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeSQLEntity(tt.args.s); got != tt.want {
				t.Errorf("escapeSQLEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}
