package customer

import (
  "fmt"
  _ "github.com/mattn/go-sqlite3"
  "log"
  "github.com/jmoiron/sqlx"
)
var db *sqlx.DB

type Customer struct {

  CustID  int    //`db:"CustID"`
  Name    string //`db:"Name"`
  Address string //`db:"Address"`
  Email   string //`db:"Email"`

}

type CustomerNotFound struct {
  custID int
}

func (e *CustomerNotFound) Error() string {
  return fmt.Sprintf("Could not find customer with CustID=%d.", e.custID)
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
    if c.CustID == 0 {
      return nil, &CustomerNotFound{id}
    } else {
      return nil, err
    }
  }
  rows.Close()
  return c, nil
}

func CreateCustomer(c *Customer) error {
  txn, err := db.Beginx()
  if err != nil {
    return err
  }
  defer txn.Rollback()
  query := "INSERT INTO customer(custid,name,address,email) VALUES(:custid, :name, :address, :email)"

  txn.NamedExec(query, &c)
  txn.Commit()
  return nil
}

func UpdateCustomer(c *Customer) error {
  upd := "UPDATE customer SET Name=:name, Address=:address, Email=:email WHERE CustID = :custid"
  txn, err := db.Beginx()
  if err != nil {
    return err
  }
  defer txn.Rollback()
  _, err = txn.NamedExec(upd, &c)
  if err != nil {
    return err
  }
  txn.Commit()
  return nil

}

func DeleteCustomer(custID int) error {
  c := Customer{CustID: custID}
  del := "DELETE FROM customer WHERE CustID=:custid"
  txn, err := db.Beginx()
  if err != nil {
    return err
  }
  defer txn.Rollback()
  _, err = txn.NamedExec(del, &c)
  if err != nil {
    return err
  }
  txn.Commit()
  //c = &Customer{}
  return nil
}

func DBClose() {
  db.Close()
}
