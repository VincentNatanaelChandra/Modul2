package Model

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ResultProduct struct {
	ID    int
	Name  string
	Price int
}

type Transaction struct {
	ID        int `json:"status"`
	UserID    int `json:"user"`
	ProductID int `json:"product"`
	Quantity  int `json:"quantity"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

type ProductResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    Product `json:"data"`
}

type ProductsResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type TransactionResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    Transaction `json:"data"`
}

type TransactionsResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"data"`
}

type TransactionProductUser struct {
	TransactionsDetail TransactionsDetail `json:"transactiondetail"`
	User               User               `json:"user"`
	Product            Product            `json:"product"`
}

type TransactionsDetailResponse struct {
	Status  int                        `json:"status"`
	Message string                     `json:"message"`
	Data    TransactionDetailResponses `json:"Data"`
}

type TransactionsDetail struct {
	ID       int     `json:"ID"`
	User     User    `json:"user"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type TransactionDetailResponses struct {
	Transaction []TransactionsDetail `json:"transaction"`
}
