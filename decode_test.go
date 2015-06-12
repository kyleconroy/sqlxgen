package dba

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type Book struct {
	ID    int    `dba:"id"`
	Shelf int    `dba:"-"`
	Title string `dba:"name"`
}

func TestUnmarshal(t *testing.T) {
	assert := assert.New(t)

	os.Remove("./test.db")
	defer os.Remove("./test.db")

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table books (id integer not null primary key, name text);
	insert into books (name) values ('GopherTales');
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM books WHERE name = $1", "GopherTales")
	if err != nil {
		t.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var book Book
		err = Unmarshal(rows, &book)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(1, book.ID)
		assert.Equal("GopherTales", book.Title)
	}

	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

}
