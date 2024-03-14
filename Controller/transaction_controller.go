package Controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	m "Modul2/Model"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM transactions"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	var transaction m.Transaction
	var transactions []m.Transaction
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}
	UserId, _ := strconv.Atoi(r.Form.Get("userid"))
	ProductId, _ := strconv.Atoi(r.Form.Get("productid"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	_, err = db.Exec("INSERT INTO transactions(userid, productid, quantity) values (?,?,?)",
		UserId,
		ProductId,
		quantity,
	)

	_, err = db.Exec("INSERT INTO products (id, name, price) VALUES (?, '', 0)", ProductId)
	if err != nil {
		log.Println("Error")
		return
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", ProductId).Scan(&count)
	if err != nil {
		log.Println("Error: No products")
		return
	}

	if count == 0 {
		_, errQuery := db.Exec("INSERT INTO products (id, name, price) VALUES (?, '', 0)", ProductId)
		if errQuery != nil {
			log.Println("Error: insertting products")
			return
		}
	}
	sendSuccessResponse(w)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	transID := r.URL.Query().Get("id")

	if transID == "" {
		log.Println("Error: ID missing")
		http.Error(w, "Bad Request: ID missing", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("userid"))
	prodID, _ := strconv.Atoi(r.Form.Get("productid"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	if userID == 0 || prodID == 0 || quantity == 0 {
		log.Println("Error: Incomplete data provided")
		http.Error(w, "Bad Request: Incomplete data", http.StatusBadRequest)
		return
	}

	data, err := db.Begin()
	if err != nil {
		log.Println("Error database not found:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("UPDATE products SET userid = ?, productid = ?, quantity = ? WHERE id = ?", userID, prodID, quantity, transID)
	if errQuery != nil {
		http.Error(w, "Update failed", http.StatusBadRequest)
		return
	}

	if errQuery == nil {
		sendSuccessResponse(w)
	} else {
		sendErrorResponse(w)
	}

}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	TransId := r.URL.Query().Get("id")

	_, errQuery := db.Exec("DELETE FROM transactions WHERE id=?",
		TransId,
	)

	if errQuery == nil {
		sendSuccessResponseTransaction(w)
	} else {
		sendErrorResponseTransaction(w)
	}

}

func sendSuccessResponseTransaction(w http.ResponseWriter) {
	var response m.TransactionResponse
	response.Status = 200
	response.Message = "Success"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponseTransaction(w http.ResponseWriter) {
	var response m.TransactionResponse
	response.Status = 400
	response.Message = "Failed"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDetailUserTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := `SELECT t.id, u.id, u.name, u.age, u.address, p.id, p.name, p.price, t.quantity
			  FROM transactions t JOIN users u ON t.id = u.id
			  JOIN products p ON t.id = p.id`
	transactionRow, err := db.Query(query)
	if err != nil {
		print(err.Error())
		sendErrorResponse(w)
		return
	}
	var transactionUser m.TransactionsDetail
	var transactionUsers []m.TransactionsDetail
	for transactionRow.Next() {
		if err := transactionRow.Scan(
			&transactionUser.ID, &transactionUser.User.ID, &transactionUser.User.Name,
			&transactionUser.User.Age, &transactionUser.User.Address, &transactionUser.Product.ID,
			&transactionUser.Product.Name, &transactionUser.Product.Price, &transactionUser.Quantity); err != nil {
			print(err.Error())
			sendErrorResponse(w)
			return
		} else {
			transactionUsers = append(transactionUsers, transactionUser)
		}
	}

	var response m.TransactionsDetailResponse
	response.Status = 200
	response.Message = "Success"
	response.Data.Transaction = transactionUsers
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDetailUserTransactionbyID(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userId := r.URL.Query().Get("id")
	query := `SELECT t.id, u.id, u.name, u.age, u.address, p.id, p.name, p.price, t.quantity
			  FROM transactions t JOIN users u ON t.id = u.id
			  JOIN products p ON t.id = p.id WHERE u.ID = ?`
	transactionRow, err := db.Query(query, userId)
	if err != nil {
		print(err.Error())
		sendErrorResponse(w)
		return
	}
	var transactionUser m.TransactionsDetail
	var transactionUsers []m.TransactionsDetail
	for transactionRow.Next() {
		if err := transactionRow.Scan(
			&transactionUser.ID, &transactionUser.User.ID, &transactionUser.User.Name,
			&transactionUser.User.Age, &transactionUser.User.Address, &transactionUser.Product.ID,
			&transactionUser.Product.Name, &transactionUser.Product.Price, &transactionUser.Quantity); err != nil {
			print(err.Error())
			sendErrorResponse(w)
			return
		} else {
			transactionUsers = append(transactionUsers, transactionUser)
		}
	}

	var response m.TransactionsDetailResponse
	response.Status = 200
	response.Message = "Success"
	response.Data.Transaction = transactionUsers
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
