package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/moio/mgr-dump/schemareader"
)

// cd spacewalk/java; make -f Makefile.docker dockerrun_pg
const connectionString = "user='spacewalk' password='spacewalk' dbname='susemanager' host='localhost' port='5432' sslmode=disable"

// psql --host=localhost --port=5432 --username=spacewalk susemanager

// go run . | dot -Tx11
func main() {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	tables := schemareader.ReadTables(db)

	fmt.Printf("graph schema {\n")
	fmt.Printf("  layout=circo;")
	fmt.Printf("  mindist=0.3;")

	for _, table := range tables {
		fmt.Printf("\"%s\" [shape=box];\n", table.Name)

		for _, column := range table.Columns {
			fmt.Printf("\"%s-%s\" [label=\"\" xlabel=\"%s\"];\n", table.Name, column, column)
			fmt.Printf("\"%s\" -- \"%s-%s\";\n", table.Name, table.Name, column)
		}
	}

	fmt.Printf("}")
}
