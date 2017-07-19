package main

import (
	"os"
	"database/sql"
	"log"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if _, err := os.Stat("./foo.db"); err != nil {
		log.Fatal(err)
	} else if os.IsNotExist(err) {
		os.Create("./foo.db")
	} else {
		os.Remove("./foo.db")
	}

	db, err := sql.Open("sqlite3", "./foo.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("Hello World %03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer  rows.Close()

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
