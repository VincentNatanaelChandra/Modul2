package Controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	m "Modul2/Model"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM products"
	name := r.URL.Query()["name"]
	price := r.URL.Query()["price"]

	if name != nil {
		fmt.Println(name[0])
		query += "WHERE name='" + name[0] + "'"
	}

	if price != nil {
		if name[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " price ='" + price[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var product m.Product
	var products []m.Product
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Println(err)
			return
		} else {
			products = append(products, product)
		}
	}

	if !rows.Next() {
		w.Header().Set("Content-Type", "application/json")
		var response m.ProductsResponse
		response.Status = 404
		response.Message = "Data not found"
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.ProductsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = products
	json.NewEncoder(w).Encode(response)
}

func GetAllProductsGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()
	err := r.ParseForm()
	if err != nil {
		sendErrorResponseProduct(w, "Error: error parsing data")
		return
	}

	var products []m.Product
	result := db.Find(&products)

	if result.Error == nil {
		sendSuccessResponseProduct(w, "Success")
	} else {
		sendErrorResponseProduct(w, "Failed")
	}
}

func InsertProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}
	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	_, errQuery := db.Exec("INSERT INTO products(name, price) values (?,?)",
		name,
		price,
	)

	if errQuery == nil {
		sendSuccessResponseProduct(w, "Success")
	} else {
		sendErrorResponseProduct(w, "Failed")
	}
}

func InsertProductGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponseProduct(w, "Error: error parsing data")
		return
	}

	name := r.Form.Get("name")
	price, err := strconv.Atoi(r.Form.Get("price"))
	if name == "" || price == 0 {
		sendErrorResponseProduct(w, "Null")
		return
	}
	if err != nil {
		sendErrorResponseProduct(w, "Invalid price")
		return
	}

	product := m.Product{Name: name, Price: price}
	result := db.Create(&product)

	if result.Error == nil {
		sendSuccessResponseProduct(w, "Success")
	} else {
		sendErrorResponseProduct(w, "Failed")
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	prodID := r.URL.Query().Get("id")

	if prodID == "" {
		log.Println("Error: ID missing")
		http.Error(w, "Bad Request: ID missing", http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	if name == "" || price == 0 {
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

	_, errQuery := db.Exec("UPDATE products SET name = ?, price = ? WHERE id = ?", name, price, prodID)
	if errQuery != nil {
		http.Error(w, "Update failed", http.StatusBadRequest)
		return
	}

	if errQuery == nil {
		sendSuccessResponseProduct(w, "Success")
	} else {
		sendErrorResponseProduct(w, "Failed")
	}

}

func UpdateProductGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	err := r.ParseForm()
	if err != nil {
		return
	}

	productID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if productID == 0 {
		sendErrorResponseProduct(w, "Bad request: Missing ID input")
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	product := m.Product{Name: name, Price: price}
	result := db.Save(&product)

	if result.Error == nil {
		sendSuccessResponseProduct(w, "Success")
	} else {
		sendErrorResponseProduct(w, "Failed")
	}
}

func UpdateProductGormRaw(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponseProduct(w, "Failed to parse form data")
		return
	}

	productID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || productID == 0 {
		sendErrorResponseProduct(w, "Bad request: Missing or invalid ID input")
		return
	}

	name := r.Form.Get("name")
	priceStr := r.Form.Get("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		sendErrorResponseProduct(w, "Invalid price")
		return
	}

	var result m.ResultProduct
	db.Raw("UPDATE products SET name = ?, price = ? WHERE id = ?", name, price, productID).Scan(&result)

	sendSuccessResponseProduct(w, "Success")
}

func DeleteSingleProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	prodID := r.URL.Query().Get("id")

	if prodID == "" {
		sendErrorResponseProduct(w, "Missing ID")
		return
	} else {
		data, err := db.Begin()
		if err != nil {
			sendErrorResponseProduct(w, "Db not found")
			return
		}
		defer data.Rollback()

		_, err = db.Exec("DELETE FROM transactions WHERE id = ?", prodID)
		if err != nil {
			sendErrorResponseProduct(w, "Delete Failed")
			return
		}

		_, err = db.Exec("DELETE FROM products WHERE id = ?", prodID)
		if err != nil {
			sendErrorResponseProduct(w, "Delete Failed")
			return
		}

		if err == nil {
			sendSuccessResponseProduct(w, "Success")
		} else {
			sendErrorResponseProduct(w, "Failed")
		}
	}
}

func DeleteProductGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	err := r.ParseForm()
	if err != nil {
		return
	}

	getID := r.URL.Query().Get("id")
	if getID == "" {
		sendErrorResponseProduct(w, "Missing ID")
		return
	}

	productID, err := strconv.Atoi(getID)
	if err != nil {
		sendErrorResponseProduct(w, "Invalid ID")
		return
	}

	product := m.Product{ID: productID}
	result := db.Delete(product)

	if result.Error == nil {
		sendSuccessResponseProduct(w, "Success")
	} else {
		sendErrorResponseProduct(w, "Failed")
	}
}

func sendSuccessResponseProduct(w http.ResponseWriter, message string) {
	var response m.ProductResponse
	response.Status = 200
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponseProduct(w http.ResponseWriter, message string) {
	var response m.ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
