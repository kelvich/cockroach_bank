package main

import (
	"fmt"
	"sync"
	"math/rand"
	"bytes"
	"database/sql"
	"time"
	_ "github.com/cockroachdb/cockroach/sql"
)

var connstrs = [...]string{
	"https://root@astro10:26257?certs=certs",
	"https://root@astro9:26257?certs=certs",
	"https://root@astro8:26257?certs=certs",
	"https://root@astro6:26257?certs=certs",
	"https://root@astro5:26257?certs=certs",
	"https://root@astro4:26257?certs=certs",
}
var numAccounts = 1000
var maxTransfer = 100
var numTransactions = 100
var numTransferWorkers = 6
var numInspectWorkers = 1

func prepare(){
	db, err := sql.Open("cockroach", connstrs[0])
	if err != nil {
		checkErr(err)
	}
	defer db.Close()

	schema := `
		CREATE TABLE IF NOT EXISTS bank.accounts (
		  id INT PRIMARY KEY,
		  balance INT NOT NULL
		)`
	if _, err = db.Exec(schema); err != nil {
		checkErr(err)
	}
	if _, err = db.Exec("TRUNCATE TABLE bank.accounts"); err != nil {
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
	stmt := `INSERT INTO bank.accounts (id, balance) VALUES ` + placeholders.String()
	if _, err = db.Exec(stmt, values...); err != nil {
		checkErr(err)
	}

}

func transfer(id int, wg *sync.WaitGroup){
	var commitErrors = 0

	db, err := sql.Open("cockroach", connstrs[id%6])
	if err != nil {
		checkErr(err)
	}
	defer db.Close()

	for i := 0; i < numTransactions; i++ {
		// from := rand.Intn(numAccounts)
		// var to int

		// for {
		// 	to = rand.Intn(numAccounts)
		// 	if from != to {
		// 		break
		// 	}
		// }
		amount := rand.Intn(maxTransfer)

		tx, err := db.Begin()
		checkErr(err)

		// _, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SNAPSHOT")
		// checkErr(err)

		update := "UPDATE bank.accounts SET balance = balance + $1 WHERE id = $2"

		_, err = tx.Exec(update, -1*amount, 2*id+1)
		checkErr(err)

		_, err = tx.Exec(update,    amount, 2*id+2)
		checkErr(err)

		err = tx.Commit()
		if err != nil {
		    commitErrors += 1
		}
		// checkErr(err)
	}


	fmt.Printf("transfer#%d: %d errors of %d total\n", id, commitErrors, numTransactions)
	wg.Done()
}


func inspect(wg *sync.WaitGroup){
	var result int32
	var resultNew int32

	db, err := sql.Open("cockroach", connstrs[0])
	if err != nil {
		checkErr(err)
	}
	defer db.Close()

	result = 1
	for {
		inspect := "select sum(balance) from bank.accounts;"
		err := db.QueryRow(inspect).Scan(&resultNew)
		checkErr(err)
		if resultNew != result {
			fmt.Printf("Total = %d\n", resultNew)
			result = resultNew
		}
	}

}

func main() {
	start := time.Now()
	prepare()
	fmt.Printf("database prepared in %0.2f seconds\n", time.Since(start).Seconds())

	var transferWg sync.WaitGroup
	var inspectWg sync.WaitGroup

	start = time.Now()
	transferWg.Add(numTransferWorkers)
	for i := 0; i < numTransferWorkers; i++ {
	    go transfer(i, &transferWg)
	}
 
	inspectWg.Add(numInspectWorkers)
	for i := 0; i < numInspectWorkers; i++ {
	    go inspect(&inspectWg)
	}

	transferWg.Wait()

	finishTime := time.Since(start).Seconds()
	fmt.Printf("writers finished in %0.2f seconds\n", finishTime)
	fmt.Printf("TPS = %0.2f\n", float64(numTransactions*numTransferWorkers)/finishTime)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}



