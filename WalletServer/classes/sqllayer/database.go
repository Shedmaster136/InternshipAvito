package sqllayer

import (
	"database/sql"
	"fmt"
	"errors"
	_"github.com/lib/pq"
)

type Database struct {
	connectHandle *sql.DB
}

func NewDatabase() (*Database, error){
	dbConn, err := sql.Open("postgres","host=localhost port=5432 user=postgres dbname=walletmgmt sslmode=disable password=75919021")
	if err != nil{
		fmt.Println("SQL: Open DB error",err)
		return nil, err
	}
	err = dbConn.Ping()
	if err != nil{
		fmt.Println("SQL: Db Connection error", err)
		return nil, err
	}
	return &Database{connectHandle: dbConn}, nil
}

func (db Database) Stop() error{
	return db.connectHandle.Close()
}

func (db Database) CreateNewUser(newUser TableUsers) error{
	_, err := db.connectHandle.Exec("INSERT INTO users (userid, userwallet) VALUES ($1, $2)", newUser.ID, newUser.Wallet)
	if err != nil{
		fmt.Println("SQL: New user creation error\n",err)
		return err
	}
	return nil
}

func (db Database) UpdateUserWallet(updatedUser TableUsers) error{
	_, err := db.connectHandle.Exec("UPDATE users SET userwallet=$2 WHERE userid=$1", updatedUser.ID, updatedUser.Wallet)
	if err != nil{
		fmt.Println("SQL: Update user table error\n",err)
		return err
	}
	return nil
}

func (db Database) GetUserByID(userID int64) (TableUsers, error){
	users, err := db.connectHandle.Query("SELECT * FROM users WHERE userid = $1", userID)
	if err != nil{
		fmt.Println("SQL ConnectHandle Error ", err)
		return TableUsers{}, err
	}
	defer users.Close()
	if(!users.Next()){
		newError := errors.New("No such user in the database")
		fmt.Println("SQL: User does not exist\n", newError)
		return TableUsers{}, newError
	}
	var user TableUsers
	err = users.Scan(&user.ID,&user.Wallet)
	if err != nil{
		fmt.Println("SQL: ...\n",err)
		return TableUsers{}, users.Err()
	}
	return user, nil
}



//structures for getting data from/to tables
type TableUsers struct {
	ID int64
	Wallet int32
}
type TableServices struct {
	Name string
	ID int64
}
type TableOrders struct {
	ID int64
	fkUser int64
	fkService int64
	Price int32
}
type TableOrderStates struct {

}
