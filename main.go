package main

import (
  "encoding/json"
	"fmt"
	"log"
	"os"
  "strconv"
  "net/http"
  "io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
  "github.com/gorilla/mux"
  c "github.com/yensho/Ticketing/customer"
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

func home(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(struct{Message string}{"This is the running server. The home endpoint only returns this message."})
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
  custID, err := strconv.Atoi(mux.Vars(r)["id"])
  if err != nil {
    panic(err)
  }

  cust, err := c.GetCustomer(custID)
  if err != nil {
    fmt.Fprintf(w, err.Error())
  } else {
    json.NewEncoder(w).Encode(cust)
  }
}

func newCustomer(w http.ResponseWriter, r *http.Request) {
  var cust c.Customer
  requestBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    fmt.Fprintf(w, "Error: Corrupted Format.")
  }

  json.Unmarshal(requestBody, &cust)
  err = c.CreateCustomer(&cust)
  if err != nil {
    log.Fatal(err)
  }
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(cust)
}

func updCustomer(w http.ResponseWriter, r *http.Request) {
  var cust c.Customer
  custID, err := strconv.Atoi(mux.Vars(r)["id"])
  if err != nil {
    panic(err)
  }
  requestBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err)
  }

  json.Unmarshal(requestBody, &cust)
  if custID != cust.CustID {
    fmt.Fprintf(w, "The CustID in the body must match the id in the url. Please correct your request.")
  } else {
    _, err := c.GetCustomer(custID)
    if err != nil{
      fmt.Fprintf(w, err.Error())
      //json.NewEncoder(w).Encode(err.Error())
    } else {
      err = c.UpdateCustomer(&cust)
      w.WriteHeader(http.StatusOK)
      json.NewEncoder(w).Encode(cust)
    }
  }
}

func delCustomer(w http.ResponseWriter, r *http.Request) {
  custID, err := strconv.Atoi(mux.Vars(r)["id"])
  if err != nil {
    panic(err)
  }

  _, err = c.GetCustomer(custID)
  if err != nil{
    fmt.Fprintf(w, err.Error())
    //json.NewEncoder(w).Encode(err.Error())
  } else {
    err = c.DeleteCustomer(custID)
    w.WriteHeader(http.StatusOK)
  }
}

func main() {
  /*
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
  */
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", home)
  router.HandleFunc("/customer", newCustomer).Methods("POST")
  router.HandleFunc("/customer/{id}", getCustomer).Methods("GET")
  router.HandleFunc("/customer/{id}", updCustomer).Methods("PUT")
  router.HandleFunc("/customer/{id}", delCustomer).Methods("DELETE")
  fmt.Println("Server now listening...")
  log.Fatal(http.ListenAndServe(":8080", router))

  c.DBClose()
}
