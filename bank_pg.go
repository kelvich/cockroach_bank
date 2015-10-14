package main

import (
	"fmt"
	"math/rand"
	"bytes"
	"database/sql"
	_ "github.com/lib/pq"
)

func prepare(){

}

func main() {
	var numAccounts = 1000
	var maxTransfer = 100
	var numTransactions = 1000


	db, err := sql.Open("postgres", "dbname=postgres sslmode=disable")
	if err != nil {
		checkErr(err)
	}
	defer db.Close()

	schema := `
		CREATE TABLE IF NOT EXISTS accounts (
		  id INT PRIMARY KEY,
		  balance INT NOT NULL
		)`
	if _, err = db.Exec(schema); err != nil {
		checkErr(err)
	}
	if _, err = db.Exec("TRUNCATE TABLE accounts"); err != nil {
		checkErr(err)
	}

	var placeholders bytes.Buffer
	var values []interface{}
	for i := 0; i < numAccounts; i++ {
		if i > 0 {
			placeholders.WriteString(", ")
		}
		fmt.Fprintf(&placeholders, "($%d, 0)", i+1)
		values = append(values, i)
	}
	stmt := `INSERT INTO accounts (id, balance) VALUES ` + placeholders.String()
	if _, err = db.Exec(stmt, values...); err != nil {
		checkErr(err)
	}


	for i := 0; i < numTransactions; i++ {
		from := rand.Intn(numAccounts)
		var to int
		for {
			to = rand.Intn(numAccounts)
			if from != to {
				break
			}
		}

		amount := rand.Intn(maxTransfer)

		tx, err := db.Begin()
		if err != nil {
		    checkErr(err)
		}

		update := "UPDATE accounts SET balance = balance + $1 WHERE id = $2"

		if _, err = tx.Exec(update, -1*amount, from); err != nil {
			checkErr(err)
		}
		if _, err = tx.Exec(update, amount, to); err != nil {
			checkErr(err)
		}
		if err = tx.Commit(); err != nil {
			checkErr(err)
		}
	}

}


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}