package main

import (
	"fmt"
	"log"
	c "main/customer"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func init() {
	os.Remove("./test.db")

	db, err := sqlx.Connect("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	c.SetDB(db)

	sqlStmt := `
  create table customer (CustID integer not null primary key, Name text, Address text, Email text);
  delete from customer;
  `

	db.MustExec(sqlStmt)
	/*
	  tx := db.MustBegin()
	  stmt, err := tx.Preparex("insert into customer(id, name, address) values (?, ?, ?)")
	  if err != nil {
	    log.Fatal(err)
	  }
	  defer stmt.Close()
	*/
	for i := 1; i < 10; i++ {
		cust := c.Customer{i, fmt.Sprint("Tester Name", i), "123 Test Rd", fmt.Sprint("Test", i, "@email.com")}
		err = c.CreateCustomer(&cust)
		if err != nil {
			log.Fatal(err)
		}
	}
	//tx.Commit()

}

func main() {

	for i := 1; i < 10; i++ {
    cust, err := c.GetCustomer(i)
    if err != nil {
      panic(err)
    }
		fmt.Printf("%#v\n", cust)
	}

  c.DBClose()
}
