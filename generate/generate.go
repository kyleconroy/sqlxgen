package generate

import (
	"database/sql"
	"fmt"
	"go/format"
	"strings"
)

func Struct(db *sql.DB, tableName, pkgName, structName string) ([]byte, error) {
	columns, err := columnInformation(db, tableName)
	if err != nil {
		return []byte{}, err
	}

	src := fmt.Sprintf("package %s\ntype %s %s\n}",
		pkgName,
		structName,
		generateFields(columns))
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

type column struct {
	Name string
	Type string
}

func columnInformation(db *sql.DB, table string) ([]column, error) {
	columns := []column{}
	rows, err := db.Query(`
		select column_name, data_type from information_schema.columns 
		where table_name = $1 order by column_name`, table)
	if err != nil {
		return columns, err
	}
	defer rows.Close()
	for rows.Next() {
		var c column
		err := rows.Scan(&c.Name, &c.Type)
		if err != nil {
			return columns, err
		}
		columns = append(columns, c)
	}
	return columns, rows.Err()
}

var postgresDataTypes = map[string]string{
	"boolean":           "bool",
	"string":            "string",
	"double precision":  "float64",
	"integer":           "int",
	"character varying": "string",
	"text":              "string",
	"date":              "time.Time",
	"timestamp with time zone":    "time.Time",
	"timestamp":                   "time.Time",
	"timestamp without time zone": "time.Time",
}

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

func generateFields(columns []column) string {
	structure := "struct {"

	for _, c := range columns {
		typeName := postgresDataTypes[c.Type]
		if typeName == "" {
			typeName = "interface{}"
		}

		parts := strings.Split(c.Name, "_")
		for i, p := range parts {
			if commonInitialisms[strings.ToUpper(p)] {
				parts[i] = strings.ToUpper(p)
			} else {
				parts[i] = strings.Title(p)
			}
		}

		structure += fmt.Sprintf("\n%s %s `db:\"%s\"`",
			strings.Join(parts, ""),
			typeName,
			c.Name)
	}
	return structure
}
