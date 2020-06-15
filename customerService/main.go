package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Customer struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	State    string `json:"state,omitempty"`
	City     string `json:"city,omitempty"`
}

var (
	customer   []Customer
	address    []Address
	dbInstance *sql.DB
)

// retrieving a single customer
func getACustomer(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	dbInstance := dbConn()
	selDB, err := dbInstance.Query("SELECT * FROM goCustomer.customer c inner join goCustomer.address a where c.id = a.customer_id and c.id= ?", params["id"])

	if err != nil {
		panic(err.Error())
	}
	cust := Customer{}

	for selDB.Next() {
		var id, firstname, lastname, address_id, address1, address2, state, city, customer_id string
		err = selDB.Scan(&id, &firstname, &lastname, &address_id, &address1, &address2, &state, &city, &customer_id)
		if err != nil {
			panic(err.Error())
		}
		cust.ID = id
		cust.Firstname = firstname
		cust.Lastname = lastname
		cust.Address = &Address{Address1: address1, Address2: address2, City: city, State: state}

		json.NewEncoder(w).Encode(cust)
	}
}

// retrieving all the customers
func getAllCustomers(w http.ResponseWriter, req *http.Request) {
	dbInstance := dbConn()
	selDB, err := dbInstance.Query("SELECT * FROM goCustomer.customer c inner join goCustomer.address a where c.id = a.customer_id ORDER BY c.id ASC")

	if err != nil {
		panic(err.Error())
	}
	cust := Customer{}

	for selDB.Next() {
		var id, firstname, lastname, address_id, address1, address2, state, city, customer_id string
		err = selDB.Scan(&id, &firstname, &lastname, &address_id, &address1, &address2, &state, &city, &customer_id)
		if err != nil {
			panic(err.Error())
		}
		cust.ID = id
		cust.Firstname = firstname
		cust.Lastname = lastname
		cust.Address = &Address{Address1: address1, Address2: address2, City: city, State: state}

		customer = append(customer, cust)
	}
	json.NewEncoder(w).Encode(customer)
}

// for deleting a customer
func deleteACustomer(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if req.Method == "DELETE" {
		db := dbConn()
		cust := params["id"]
		fmt.Println(cust)
		delForm, err := db.Prepare("DELETE FROM customer WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(cust)
		fmt.Println("Customer deleted")
	}
}

// for creating a customer
func createACustomer(w http.ResponseWriter, req *http.Request) {
	db := dbConn()
	if req.Method == "POST" {
		var customer Customer
		_ = json.NewDecoder(req.Body).Decode(&customer)
		json.NewEncoder(w).Encode(customer)

		insForm, err := db.Prepare("INSERT INTO `customer` (`firstname`, `lastname`) VALUES (?,?);")
		if err != nil {
			panic(err.Error())
		}
		res, err := insForm.Exec(customer.Firstname, customer.Lastname)
		id, _ := res.LastInsertId()

		insForm, err2 := db.Prepare("INSERT INTO `address` (`address1`, `address2`, `state`, `city`, `customer_id`) VALUES (?, ?, ?, ?, ?)")
		if err2 != nil {
			panic(err.Error())
		}
		insForm.Exec(customer.Address.Address1, customer.Address.Address2, customer.Address.State, customer.Address.City, id)
	}
}

// for updating a customer
func updateACustomer(w http.ResponseWriter, req *http.Request) {
	db := dbConn()
	if req.Method == "PUT" {
		var customer Customer
		_ = json.NewDecoder(req.Body).Decode(&customer)
		json.NewEncoder(w).Encode(customer)

		params := mux.Vars(req)
		id := params["id"]

		updateForm, err := db.Prepare("UPDATE `customer` SET `firstname` = '?', `lastname` = '?' WHERE `customer`.`id` = ?")
		if err != nil {
			panic(err.Error())
		}
		updateForm.Exec(customer.Firstname, customer.Lastname, id)

		updateForm, err2 := db.Prepare("UPDATE `address` SET `address1` = '?', `address2` = '?', `state` = '?', `city` = '?' WHERE `address`.`customer_id` = ?;")
		if err2 != nil {
			panic(err.Error())
		}
		updateForm.Exec(customer.Address.Address1, customer.Address.Address2, customer.Address.State, customer.Address.City, id)
	}
}

func main() {
	router := mux.NewRouter()
	router.Use(addJSONHeaders)

	//allowed routes
	router.HandleFunc("/customerService", getAllCustomers).Methods("GET")
	router.HandleFunc("/customerService/{id}", getACustomer).Methods("GET")
	router.HandleFunc("/customerService", createACustomer).Methods("POST")
	router.HandleFunc("/customerService/{id}", updateACustomer).Methods("PUT")
	router.HandleFunc("/customerService/{id}", deleteACustomer).Methods("DELETE")

	// routes allowed with port 3037
	log.Fatal(http.ListenAndServe(":3037", router))
	defer dbInstance.Close()
}
func addJSONHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func dbConn() (db *sql.DB) {
	if dbInstance == nil {
		db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/goCustomer")
		if err != nil {
			panic(err.Error())
		} else {

		}
		dbInstance = db

	}
	return dbInstance
}
