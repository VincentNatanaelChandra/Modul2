package main

import (
	"Modul2/Controller"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", Controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/users", Controller.InsertUser).Methods("POST")
	router.HandleFunc("/users", Controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/login", Controller.DeleteUser).Methods("POST")

	router.HandleFunc("/products", Controller.GetAllProducts).Methods("GET")
	router.HandleFunc("/v2/products", Controller.GetAllProductsGorm).Methods("GET")
	router.HandleFunc("/products", Controller.InsertProduct).Methods("POST")
	router.HandleFunc("/v2/products", Controller.InsertProductGorm).Methods("POST")
	router.HandleFunc("/v2/products", Controller.UpdateProductGormRaw).Methods("PUT")
	router.HandleFunc("/products", Controller.DeleteSingleProduct).Methods("DELETE")
	router.HandleFunc("/v2/products", Controller.DeleteProductGorm).Methods("DELETE")

	router.HandleFunc("/transactions", Controller.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions", Controller.InsertTransaction).Methods("POST")
	router.HandleFunc("/transactions", Controller.DeleteTransaction).Methods("DELETE")

	router.HandleFunc("/transactionsDetail", Controller.GetDetailUserTransaction).Methods("GET")
	router.HandleFunc("/transactionsDetailUser", Controller.GetDetailUserTransactionbyID).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
