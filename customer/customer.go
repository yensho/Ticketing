package customer

import (
  "fmt"
  _ "github.com/mattn/go-sqlite3"
  "log"
  "github.com/jmoiron/sqlx"
)
var db *sqlx.DB

type Customer struct {
  
  CustID  int    `db:"CustID"`
  Name    string `db:"Name"`
  Address string `db:"Address"`
  Email   string `db:"Email"`

}

func Connect(path string) {
  var err error
  db, err = sqlx.Open("sqlite3", path)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Connected to "+ path)
}

func SetDB(newdb *sqlx.DB) {
  db = newdb
}

func GetCustomer(id int) (*Customer, error) {
  
  rows, err := db.Queryx("SELECT CustID, Name, Address, Email FROM customer WHERE CustID = ?", id)
  if err != nil {
    return nil, err
  }

  c := &Customer{}
  rows.Next()
  if err := rows.StructScan(c); err != nil {  
    fmt.Printf("%#v\n", *c)
    return nil, err
  } 
  rows.Close()
  return c, nil
}

func CreateCustomer(c *Customer) error {
  txn, err := db.Begin()
  if err != nil {
    return err
  }
  defer txn.Rollback()
  query := "INSERT INTO customer(CustID,Name,Address,Email) VALUES(?, ?, ?, ?)"
  stmt, err := txn.Prepare(query)
  if err != nil {
    return err
  }

  _, err = stmt.Exec(c.CustID, c.Name, c.Address, c.Email)
  if err != nil  {
    return err
  }
  txn.Commit()
  return nil
}


func (c *Customer) GetID() int {
  return c.CustID
}

func (c *Customer) GetName() string {
  return c.Name
}

func (c *Customer) GetAddress() string {
  return c.Address
}

func (c *Customer) GetEmail() string {
  return c.Email
}

func DBClose() {
  db.Close()
}
