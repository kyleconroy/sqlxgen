package dal

import (
	"errors"
	"reflect"
)

func scanValues(v reflect.Value, t *typeInfo, flags fieldFlags) []interface{} {
	elem := v.Elem()
	values := []interface{}{}
	for _, field := range t.fields {
		if field.flags&flags == 1 {
			values = append(values, elem.FieldByIndex(field.idx).Addr().Interface())
		}
	}
	return values
}

func tableColumns(t *typeInfo, flags fieldFlags) (string, []string) {
	columns := []string{}
	table := t.daltable.name
	for _, field := range t.fields {
		if field.flags&flags == 1 {
			if flags == QueryInsert {
				columns = append(columns, field.name)
			} else {
				columns = append(columns, table+"."+field.name)
			}
		}
	}
	return table, columns
}

func Values(v interface{}, query fieldFlags) ([]interface{}, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return []interface{}{}, errors.New("non-pointer passed to dal.Values")
	}
	tinfo, err := getTypeInfo(val.Elem().Type())
	if err != nil {
		return []interface{}{}, err
	}
	return scanValues(val, tinfo, query), nil
}

func Expression(v interface{}, query fieldFlags) (string, []string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return "", []string{}, errors.New("non-pointer passed to dal.Values")
	}
	tinfo, err := getTypeInfo(val.Elem().Type())
	if err != nil {
		return "", []string{}, err
	}
	table, cols := tableColumns(tinfo, query)
	return table, cols, nil
}

func ExpressionRow(v interface{}, query fieldFlags) (string, []string, []interface{}, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return "", []string{}, []interface{}{}, errors.New("non-pointer passed to dal.Values")
	}
	tinfo, err := getTypeInfo(val.Elem().Type())
	if err != nil {
		return "", []string{}, []interface{}{}, err
	}
	table, cols := tableColumns(tinfo, query)
	return table, cols, scanValues(val, tinfo, query), nil
}
