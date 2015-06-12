package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/kyleconroy/dba"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	ID    int    `dba:"id"`
	Shelf int    `dba:"-"`
	Title string `dba:"name"`
}

func main() {
	os.Remove("./test.db")
	defer os.Remove("./test.db")

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table books (id integer not null primary key, name text);
	insert into books (name) values ('Gophers of Wrath');
	insert into books (name) values ('A Tale of Two Gophers');
	insert into books (name) values ('Gopherbumps');
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var book Book
		err = dba.Unmarshal(rows, &book)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(book.ID, book.Title)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}
