# Go-AppWithCrudOperations
App generated in GO for CRUD operations

"customerService" is a small APP service written in GO language to perform CRUD operations. This enables 5 API endpoints to GET, POST, PUT and DELETE a customer(s). These can be checked in the POSTMAN.

1. router.HandleFunc("/customerService", getAllCustomers).Methods("GET")
2. router.HandleFunc("/customerService/{id}", getACustomer).Methods("GET")
3. router.HandleFunc("/customerService", createACustomer).Methods("POST")
4. router.HandleFunc("/customerService/{id}", updateACustomer).Methods("PUT")
5. router.HandleFunc("/customerService/{id}", deleteACustomer).Methods("DELETE")

  
"customerService" directory consists of two files main.go and DB.sql

This service listens to port "3037" and my Mysql DB runs on port "8889".
