package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/kyleconroy/sqlxgen/generate"
	_ "github.com/lib/pq"
)

func main() {
	var pkg = flag.String("package", "main", "Package name, defaults to main")
	var name = flag.String("struct", "<table>", "Struct name, defaults to table name")

	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Println("sqlxgen [-package] [-struct] <database-url> <table>")
		flag.PrintDefaults()
		return
	}

	url := flag.Arg(0)
	table := flag.Arg(1)

	if *name == "<table>" {
		*name = table
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	blob, err := generate.Struct(db, table, *pkg, *name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(blob))
}
