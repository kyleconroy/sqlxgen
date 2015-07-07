# sqlxgen

Generate SQLx tagged structs from existing database tables

## Installation

This package can be installed with the go get command:

```
go get github.com/kyleconroy/sqlxgen
```

## Usage

```
$ sqlxgen -pkg=barnes -struct=Book 'dbname=noble' books
package barnes

type Book struct {
  Author string   `db:"author"`
  ISBN   string   `db:"isbn"`
  Price  int      `db:"price"`
}
```

## Documentation

API documentation can be found here: http://godoc.org/github.com/kyleconroy/sqlxgen
