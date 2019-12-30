package main

import (
	"fmt"
	"log"
	c "main/customer"
	"os"

  //"net/http"

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
  create table customer (custid integer not null primary key, name text, address text, email text);
  delete from customer;
  `

	db.MustExec(sqlStmt)
	for i := 1; i < 10; i++ {
		cust := c.Customer{i, fmt.Sprint("Tester Name", i), "123 Test Rd", fmt.Sprint("Test", i, "@email.com")}
		err = c.CreateCustomer(&cust)
		if err != nil {
			log.Fatal(err)
		}
	}


}

func main() {

  mycust := &c.Customer{CustID: 5, Name: "Bob Test", Address: "54321 Test Dr", Email: "Bob@Test.com"}
  c.UpdateCustomer(mycust)
  
	for i := 1; i < 10; i++ {
    cust, err := c.GetCustomer(i)
    if err != nil {
      panic(err)
    }
		fmt.Printf("%#v\n", cust)
	}
  cust, _ := c.GetCustomer(4)
  c.DeleteCustomer(cust)
  cust, err := c.GetCustomer(4)
  if err != nil {
    fmt.Println(err)
  }
  c.DBClose()
}
