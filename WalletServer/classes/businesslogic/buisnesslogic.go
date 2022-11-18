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
		err = bl.appDatabase.CreateNewUser(sqllayer.TableUsers{ID: userID.ID, Wallet: userID.TopUp})
		if err != nil {
			fmt.Println("Create err")
			return err
		}
	}else{
		existingUser.Wallet += userID.TopUp
		err = bl.appDatabase.UpdateUserWallet(*existingUser)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bl BusinessLogic) GetUserBalance(userID int64)  (int32, error){
	existingUser, err := bl.appDatabase.GetUserByID(userID)
	if err != nil{
		return 0, err
	}else{
		return existingUser.Wallet, err
	}
	
}

type User struct {
	ID    int64 "json:\"id\""
	TopUp int32 "json:\"topup\""
}
