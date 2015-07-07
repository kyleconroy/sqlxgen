# Database Awesome

## Description

Easily unmarshal database rows into structs.

```go
package main

import (
	"database/sql"
	"log"

	"github.com/kyleconroy/dba"
)

type Book struct {
	ID    int    `dba:"id"`
	Shelf int    `dba:"-"`
	Title string `dba:"name"`
}

func main() {
	db, _ := sql.Open("...", "...")
	rows, _ := db.Query("SELECT * FROM books")
	rows.Next()

	var book Book
	dba.Unmarshal(rows, &book)
	log.Println(book.ID, book.Title)
}
```

## Installation

This package can be installed with the go get command:

```
go get github.com/kyleconroy/dba/dba
```

## Usage

DBA can generate structs directly from database tables to save you time

```
$ dba generate --url="dbname=noble" --table=books --struct=Book --pkg=barnes
package barnes

type Book struct {
  Author string   // author
  ISBN   string   // isbn
  Price  int      // price
}
```

## Documentation

API documentation can be found here: http://godoc.org/github.com/kyleconroy/dba

Examples can be found under the ./_example directory
