package main

import (
	"WalletServer/classes/wserver"
	"WalletServer/classes/sqllayer"
	"WalletServer/classes/handlers"
	"WalletServer/classes/businesslogic"
	"bufio"
	"os"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	database, err:= sqllayer.NewDatabase()
	if err != nil{
		panic(err)
	}
	defer database.Stop()
	logic := businesslogic.NewBusinessLogic(database)
	handlers := handlers.NewHandlers(logic)
	server := new(wserver.WebServer)
	go func() {
		err := server.Start("9090", handlers)
		if err != nil{
			panic(err)
		}
		fmt.Println("server started")
	}()
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	server.Stop()
}