package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jrcichra/discord-day-trader/common"

	//mysql
	_ "github.com/go-sql-driver/mysql"
)

//Database - object to talk to the database
type Database struct {
	dbh *sql.DB
}

//Connect - connect to the database given a connect string
func (d *Database) Connect(dsn string) error {
	// connect to the database
	var err error
	d.dbh, err = sql.Open("mysql", dsn)
	if err == nil {
		err = d.dbh.Ping()
	}
	return err
}

//Reconnect - keep trying to connect to the database
func (d *Database) Reconnect(dsn string) {
	var err error
	err = nil
	for err == nil {
		err = d.Connect(dsn)
		if err != nil {
			log.Println(err)
			d.dbh.Close()
			time.Sleep(time.Duration(1) * time.Second)
			err = nil
		}

	}
}

//Database level query operations

//GetUser - returns all the user information for a given id
func (d *Database) GetUser(id common.UserID) (*common.User, error) {
	rows, err := d.dbh.Query("SELECT user_id, username, registered, last_action FROM users WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rowcount := 0
	var user common.User
	for rows.Next() {
		rowcount++
		err := rows.Scan(&user.UserID, &user.Username, &user.Registered, &user.LastAction)
		if err != nil {
			return nil, err
		}
		break
	}
	if rowcount < 1 {
		err = errors.New("Missing user")
	}
	return &user, err
}

//GetAccounts - get a list of account ids tied to a user
func (d *Database) GetAccounts(id common.UserID) []*common.Account {
	rows, err := d.dbh.Query("SELECT account_id, account_name, user_id, created, account_status_id FROM accounts WHERE account_id = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var accounts []*common.Account
	for rows.Next() {
		var acc common.Account
		err := rows.Scan(&acc.AccountID, &acc.Name, &acc.UserID, &acc.Created, &acc.AccountStatusID)
		if err != nil {
			panic(err)
		}
		accounts = append(accounts, &acc)
	}
	return accounts
}

//GetAccount - get an account object from an account ID
func (d *Database) GetAccount(id common.AccountID) *common.Account {
	rows, err := d.dbh.Query("SELECT account_id, account_name, user_id, created, account_status_id FROM accounts WHERE account_id = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var acc common.Account
	for rows.Next() {
		err := rows.Scan(&acc.AccountID, &acc.Name, &acc.UserID, &acc.Created, &acc.AccountStatusID)
		if err != nil {
			panic(err)
		}
		break
	}
	return &acc
}

//GetAccountStatuses - get all possible account statuses
func (d *Database) GetAccountStatuses() []*common.AccountStatus {
	rows, err := d.dbh.Query("SELECT account_status_id, name, type FROM account_statuses")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var as []*common.AccountStatus
	for rows.Next() {
		var a common.AccountStatus
		err := rows.Scan(&a.AccountStatusID, &a.Name, &a.Type)
		if err != nil {
			panic(err)
		}
		as = append(as, &a)
	}
	return as
}

//UpdateAccountName - update an account object from an account ID
func (d *Database) UpdateAccountName(id common.AccountID, name string) error {
	_, err := d.dbh.Exec("UPDATE accounts SET account_name = ? WHERE account_id = ?", name, id)
	return err
}

//UpdateAccountStatus - update an account status from an account ID
func (d *Database) UpdateAccountStatus(id common.AccountID, state common.AccountStatusID) error {
	_, err := d.dbh.Exec("UPDATE accounts SET account_status_id = ? WHERE account_id = ?", state, id)
	return err
}

//InsertTransaction - add the transaction to the database
func (d *Database) InsertTransaction(t common.Transaction) (common.TransactionID, error) {
	res, err := d.dbh.Exec("INSERT INTO transactions (from_symbol, to_symbol, quantity, sender, receiver) VALUES (?,?,?,?,?)", t.FromSymbol, t.ToSymbol, t.Quantity, t.Sender, t.Receiver)
	id, _ := res.LastInsertId()
	return common.TransactionID(id), err
}

//CreateUser - create a user for day-trader
func (d *Database) CreateUser(u *common.User) (common.UserID, error) {
	res, err := d.dbh.Exec("INSERT INTO users (user_id,username) VALUES (?,?)", u.UserID, u.Username)
	id, _ := res.LastInsertId()
	return common.UserID(id), err
}

//GetTransaction - get a transaction from the database
func (d *Database) GetTransaction(id common.TransactionID) *common.Transaction {
	rows, err := d.dbh.Query("SELECT transaction_id, transaction_date, from_symbol, to_symbol, quantity, sender, receiver FROM transactions WHERE transaction_id = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var t common.Transaction
	for rows.Next() {
		err := rows.Scan(&t.TransactionID, &t.Date, &t.FromSymbol, &t.ToSymbol, &t.Quantity, &t.Sender, &t.Receiver)
		if err != nil {
			panic(err)
		}
		break
	}
	return &t
}

//GetAccountTransactions - get transactions from the database for an account
func (d *Database) GetAccountTransactions(id common.AccountID) []*common.Transaction {
	rows, err := d.dbh.Query("SELECT transaction_id, transaction_date, from_symbol, to_symbol, quantity, sender, receiver FROM transactions WHERE sender = ? OR receiver = ?", id, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var ts []*common.Transaction
	for rows.Next() {
		var t common.Transaction
		err := rows.Scan(&t.TransactionID, &t.Date, &t.FromSymbol, &t.ToSymbol, &t.Quantity, &t.Sender, &t.Receiver)
		if err != nil {
			panic(err)
		}
		ts = append(ts, &t)
	}
	return ts
}
