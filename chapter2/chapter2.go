package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

func main() {
	db, err := sql.Open("mysql",
		"root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	username, err := selectName(db, 2)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Println(username)
}

func selectName(db *sql.DB, id int) (string, error) {
	var username string
	selectSql := `SELECT username FROM user WHERE id= ?`
	err := db.QueryRow(selectSql, id).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "old_datum", nil
		}
		return "", errors.Wrap(err, "select name error")
	}
	return username, nil
}
