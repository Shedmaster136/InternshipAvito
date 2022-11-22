package businesslogic

import (
	"WalletServer/classes/sqllayer"
	"fmt"
)

type BusinessLogic struct {
	appDatabase *sqllayer.Database
}

func NewBusinessLogic(sqlDB *sqllayer.Database) *BusinessLogic {
	return &BusinessLogic{appDatabase: sqlDB}
}

func (bl BusinessLogic) TopUpWallet(userID User) error {
	existingUser, err := bl.appDatabase.GetUserByID(userID.ID)
	if err != nil{
		err = bl.appDatabase.CreateNewUser(sqllayer.TableUsers{ID: userID.ID, Wallet: userID.Balance})
		if err != nil {
			fmt.Println("BL: Error registering data about a new user")
			return err
		}
	}else{
		existingUser.Wallet += userID.Balance
		err = bl.appDatabase.UpdateUserWallet(existingUser)
		if err != nil {
			fmt.Println("BL: Error putting the money on the wallet")
			return err
		}
	}
	return nil
}

func (bl BusinessLogic) GetUserBalance(userID int64)  (int32, error){
	existingUser, err := bl.appDatabase.GetUserByID(userID)
	if err != nil{
		fmt.Println("BL: Error Getting user's balance")
		return 0, err
	}
	return existingUser.Wallet, nil
}

func (bl BusinessLogic) CreateOrder(newOrder Order) (error){
	order := sqllayer.TableOrders{
		ID: newOrder.OrderID, 
		FKUser: newOrder.UserID, 
		FKService: newOrder.ServiceID,
		Price: newOrder.Price,
	}
	service := sqllayer.TableServices{
		ID: newOrder.ServiceID,
		Name: newOrder.ServiceName,
	}
	user := sqllayer.TableUsers{
		ID: newOrder.UserID,
	}
	err := bl.appDatabase.CreateOrder(user,service,order)
	if err != nil{
		fmt.Println("BL: Error Creating New Order")
		return err
	}
	return nil
}

type User struct {
	ID    int64 "json:\"uid\""
	Balance int32 "json:\"balance\""
}
type Order struct{
	UserID int64 "json:\"uid\""
	OrderID int64 "json:\"oid\""
	ServiceID int64 "json:\"sid\""
	ServiceName string "json:\"sname\""
	Price int32 "json:\"price\""
}
