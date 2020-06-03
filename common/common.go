package common

import "github.com/go-sql-driver/mysql"

//Pretty much the database schema in golang objects

//UserID -
type UserID string

//User - A user of discord-day-trader
type User struct {
	UserID     UserID
	Username   string
	Registered mysql.NullTime
	LastAction mysql.NullTime
}

//AccountStatusID -
type AccountStatusID int64

//AccountStatus - Different account statuses
type AccountStatus struct {
	AccountStatusID AccountStatusID
	Name            string
	Type            string
}

//AccountID -
type AccountID int64

//Account - An account of a discord-day-trader user
type Account struct {
	AccountID       AccountID
	Name            string
	UserID          UserID
	Created         mysql.NullTime
	AccountStatusID AccountStatusID
}

//Symbol - AAPL
type Symbol string

//TransactionID -
type TransactionID int64

//Transaction - a single transaction
type Transaction struct {
	TransactionID TransactionID
	Date          mysql.NullTime
	FromSymbol    Symbol
	ToSymbol      Symbol
	Quantity      float64
	Sender        AccountID
	Receiver      AccountID
}

//Transactions - an array of Transaction structs
type Transactions []Transaction

//OrderStatusID -
type OrderStatusID int64

//OrderStatus -
type OrderStatus struct {
	OrderStatusID OrderStatusID
	Name          string
	Type          string
}

//OrderTypeID -
type OrderTypeID int64

//OrderType -
type OrderType struct {
	OrderTypeID OrderTypeID
	name        string
	Type        string
}

//OrderID -
type OrderID int64

//Order - orders to buy stocks
type Order struct {
	OrderID     OrderID
	StatusID    OrderStatusID
	OrderTypeID OrderTypeID
	AccountID   AccountID
	LimitPrice  float64
}
