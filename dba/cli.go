package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/docopt/docopt-go"
	"github.com/kyleconroy/dba"

	_ "github.com/lib/pq"
)

func main() {
	usage := `Database Awesome

Usage:
  dba generate --url=<db> --table=<tbl> [--package=<pkg>] [--struct=<name>]
  dba -h | --help
  dba --version

Options:
  -h --help        Show this screen.
  --version        Show version.
  --url=<db>       Database URL
  --table=<tbl>    Database table
  --tags           Include DBA tags in the struct
  --struct=<name>  Struct name [default: --table]
  --package=<pkg>  Package name [default: main].`

	arguments, err := docopt.Parse(usage, nil, true, "Database Awesome 1.0", false)
	if err != nil {
		log.Fatal(err)
	}

	url := arguments["--url"].(string)
	table := arguments["--table"].(string)
	pkg := arguments["--package"].(string)
	name := arguments["--struct"].(string)

	if name == "--table" {
		name = table
	}
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	_, tags := arguments["--tags"]

	blob, err := dba.Generate(db, table, pkg, name, tags)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(blob))
}
