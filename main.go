package dal

import (
	"database/sql"
	"reflect"
)

func Open(driver, url string) (*DB, error) {
	db, err := sql.Open(driver, url)
	return &DB{db: db}, err
}

// A Table represents a table defined in a relational database.
type Table struct {
	Name, Encoding, Locale string
}

type DB struct {
	db     *sql.DB
	where  []map[string][]interface{}
	orders []string
	limit  int
}

func (db *DB) clone() *DB {
	return &DB{
		where:  db.where,
		orders: db.orders,
		limit:  db.limit,
	}
}

func (db *DB) DB() *sql.DB {
	return db.db
}

func (db *DB) Limit(l int) *DB {
	db.limit = l
	return db
}

func (db *DB) Order(o string) *DB {
	db.orders = append(db.orders, o)
	return db
}

func (db *DB) Where(expression string, args ...interface{}) *DB {
	db.where = append(db.where, map[string][]interface{}{expression: args})
	return db
}

// These functions actually perform a database query
func (db *DB) All(v interface{}) *Error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return &Error{message: "non-pointer passed to All"}
	}
	_, _ = getTypeInfo(val.Elem().Type())
	return nil
}

func (db *DB) Get(v interface{}) *Error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return &Error{message: "non-pointer passed to Get"}
	}
	_, _ = getTypeInfo(val.Elem().Type())
	return nil
}

func (db *DB) Save(v interface{}) *Error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return &Error{message: "non-struct passed to Save"}
	}
	_, _ = getTypeInfo(val.Type())
	return nil
}

type Error struct {
	message string
}

func (e *Error) Failed() bool {
	return e != nil
}

// Maybe not found?
func (e *Error) Empty() bool {
	return e != nil
}

func (e *Error) Error() string {
	return e.message
}
