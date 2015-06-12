package dba

import (
	"database/sql"
	"errors"
	"reflect"
)

func Unmarshal(rows *sql.Rows, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("non-pointer passed to dba.Unmarshal")
	}
	tinfo, err := getTypeInfo(val.Elem().Type())
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	vals := make([]interface{}, len(columns))
	lookup := map[string]int{}

	for i, c := range columns {
		lookup[c] = i
		vals[i] = new(sql.RawBytes)
	}

	elem := val.Elem()

	for _, field := range tinfo.fields {
		if i, ok := lookup[field.name]; ok {
			vals[i] = elem.FieldByIndex(field.idx).Addr().Interface()
		}
	}

	return rows.Scan(vals...)
}
