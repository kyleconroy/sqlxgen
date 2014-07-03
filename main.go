package dal

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func Open(driver, url string) (*DB, error) {
	sqldb, err := sql.Open(driver, url)
	db := DB{
		db:     sqldb,
		orders: []string{},
		driver: driver,
		where:  []map[string][]interface{}{},
	}
	return &db, err
}

// A Table represents a table defined in a relational database.
type Table struct {
	Name, Encoding, Locale string
}

type DB struct {
	db     *sql.DB
	driver string
	where  []map[string][]interface{}
	orders []string
	limit  int
}

func (db *DB) clone() *DB {
	return &DB{
		db:     db.db,
		driver: db.driver,
		where:  db.where,
		orders: db.orders,
		limit:  db.limit,
	}
}

func (db *DB) DB() *sql.DB {
	return db.db
}

func (db *DB) Limit(l int) *DB {
	d := db.clone()
	d.limit = 1
	return d
}

func (db *DB) Order(o string) *DB {
	d := db.clone()
	d.orders = append(d.orders, o)
	return d
}

func (db *DB) From(tables ...string) *DB {
	return db
}

func (db *DB) Where(expression string, args ...interface{}) *DB {
	d := db.clone()
	d.where = append(d.where, map[string][]interface{}{expression: args})
	return d
}

func wrapError(err error) *Error {
	return &Error{message: err.Error(), failure: true}
}

func scanValues(v reflect.Value, t *typeInfo) []interface{} {
	elem := v.Elem()
	values := []interface{}{}
	for _, field := range t.fields {
		values = append(values, elem.FieldByIndex(field.idx).Addr().Interface())
	}
	return values
}

// Naive SQL construction here
func (db *DB) toSQL(t *typeInfo) (string, []interface{}) {
	columns := []string{}
	args := []interface{}{}
	table := t.daltable.name
	for _, field := range t.fields {
		columns = append(columns, table+"."+field.name)
	}
	q := "SELECT " + strings.Join(columns, ", ") + " FROM " + table

	if len(db.where) > 0 {
		exprs := []string{}
		// FIXME: Three nested for loops? Has to be a better way to do this
		for _, whereMap := range db.where {
			for expr, holders := range whereMap {
				exprs = append(exprs, expr)
				for _, arg := range holders {
					args = append(args, arg)
				}
			}
		}
		q = q + " WHERE " + strings.Join(exprs, " AND ")
	}

	for _, order := range db.orders {
		q = q + " ORDER BY " + order
	}

	if db.limit > 0 {
		q = q + " LIMIT ?"
		args = append(args, db.limit)
	}
	return db.replacePlaceholders(q), args
}

func (db *DB) replacePlaceholders(sql string) string {
	if db.driver != "postgres" {
		return sql
	}
	buf := &bytes.Buffer{}
	for i := 1; ; i++ {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		buf.WriteString(sql[:p])
		fmt.Fprintf(buf, "$%d", i)
		sql = sql[p+1:]
	}
	buf.WriteString(sql)
	return buf.String()
}

// These functions actually perform a database query
func (db *DB) All(v interface{}) *Error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Slice {
		return &Error{message: "non-slice passed to All"}
	}

	typ := val.Type().Elem()

	// val is a pointer to a slice of structs
	tinfo, err := getTypeInfo(typ)
	if err != nil {
		return wrapError(err)
	}

	query, args := db.toSQL(tinfo)

	rows, err := db.db.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Error{message: "no records found", empty: true}
		} else {
			return wrapError(err)
		}
	}

	defer rows.Close()
	for rows.Next() {
		rowStruct := reflect.New(typ)
		values := scanValues(rowStruct, tinfo)
		if err := rows.Scan(values...); err != nil {
			return wrapError(err)
		}
		log.Printf("%+v", rowStruct)
		val.Set(reflect.Append(val, rowStruct.Elem()))
	}

	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return &Error{message: "no records found", empty: true}
		} else {
			return wrapError(err)
		}
	}

	return nil
}

func (db *DB) Get(v interface{}) *Error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return &Error{message: "non-pointer passed to Get"}
	}

	tinfo, err := getTypeInfo(val.Elem().Type())
	if err != nil {
		return wrapError(err)
	}

	query, args := db.Limit(1).toSQL(tinfo)
	values := scanValues(val, tinfo)

	log.Println(query)

	err = db.db.QueryRow(query, args...).Scan(values...)
	switch {
	case err == sql.ErrNoRows:
		return &Error{message: "no record found", empty: true}
	case err != nil:
		return wrapError(err)
	}
	return nil
}

func (db *DB) Save(v interface{}) *Error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return &Error{message: "non-struct passed to Save", failure: true}
	}
	_, _ = getTypeInfo(val.Type())
	return nil
}

type Error struct {
	message string
	failure bool
	empty   bool
}

func (e *Error) Failed() bool {
	if e == nil {
		return false
	}
	return e.failure
}

// Maybe not found?
func (e *Error) Empty() bool {
	if e == nil {
		return false
	}
	return e.empty
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}
