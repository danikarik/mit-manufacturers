package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db         *sqlx.DB
	connString string
	sqlHeader  = "INSERT INTO `manufacturers` (`brand_name`, `country_id`, `currency_id`, `agent`, `created_at`, `updated_at`)"
)

func main() {
	var err error

	flag.StringVar(&connString, "conn", "", "Connection string")
	flag.Parse()

	db, err = sqlx.Open("mysql", connString)
	check(err)

	err = db.Ping()
	check(err)

	var statements []string
	var rows *sql.Rows
	rows, err = db.Query(`select distinct(title) as title, country from registry_manufacturers order by title`)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var title, country string
		var countryID int
		var c *sql.Rows
		err = rows.Scan(&title, &country)
		check(err)
		c, err = db.Query(`select id from countries where name = ?`, country)
		check(err)
		if c.Next() {
			err = c.Scan(&countryID)
			check(err)
			c.Close()
		}
		if len(title) > 1 && countryID > 0 {
			data := fmt.Sprintf("('%s',%d,3,'',NOW(),NOW())", title, countryID)
			statements = append(statements, data)
		}
	}

	sqlFileName := "./data/manufacturers.sql"
	file, err := os.Create(sqlFileName)
	check(err)
	defer file.Close()

	file.WriteString(sqlHeader)
	file.Write([]byte("\n"))
	file.WriteString("VALUES")

	for i, s := range statements {
		file.Write([]byte("\n\t"))
		file.WriteString(s)
		if i == len(statements)-1 {
			file.WriteString(";")
		} else {
			file.WriteString(",")
		}
	}

	fmt.Println("Done.")
}

func clean(str string) string {
	return strings.TrimSpace(strings.Trim(str, ""))
}

func check(err error) {
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
