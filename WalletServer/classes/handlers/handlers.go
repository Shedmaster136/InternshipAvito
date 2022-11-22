package handlers

import (
	"WalletServer/classes/businesslogic"
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
)

const AcceptedContentType string = "application/json"

type VerbHandler struct {
	http.Handler
	appLogic *businesslogic.BusinessLogic
}

//Handles requests to /wallet endpoints
type WalletHandler struct {
	*VerbHandler
}

//Handles requests to /transaction endpoints
type TransactionHandler struct {
	*VerbHandler
}

//how to check a mime type?
func (wh WalletHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	//better use function if contains
	if(!strings.Contains(request.Header.Get("Content-Type"), AcceptedContentType)){
		http.Error(response, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}
	
	var err error
	defer func(){
		if err != nil {
			http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	if request.Method == http.MethodPut {
		wh.putMethod(response, request)
	}else if request.Method == http.MethodGet{
		wh.getMethod(response, request)
	}else{
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (wh WalletHandler) putMethod(response http.ResponseWriter,request *http.Request) {
	
	user := businesslogic.User{}
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(response, "Unknown JSON arguments", http.StatusUnprocessableEntity)
		return
	}
	
	err = wh.appLogic.TopUpWallet(user)
	if err != nil {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return 
	}
	json.NewEncoder(response).Encode(http.StatusOK)
	return 
}



func (wh WalletHandler) getMethod(response http.ResponseWriter,request *http.Request){
	var user businesslogic.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		fmt.Println("Get Wallet error", err)
		return 
	}
	
	user.Balance, err = wh.appLogic.GetUserBalance(user.ID)
	if err != nil {
		//log
		return 
	}
	err = json.NewEncoder(response).Encode(user)
	return 
}

func NewHandlers(appLogic *businesslogic.BusinessLogic) map[string]http.Handler {
	vh := &VerbHandler{appLogic: appLogic}
	return map[string]http.Handler{
		"/wallet": WalletHandler{
			VerbHandler: vh,
		},
		"/transaction": TransactionHandler{
			VerbHandler: vh,	
		},
	}
}


func (th TransactionHandler) putMethod(response http.ResponseWriter, request *http.Request){
	var order businesslogic.Order
	err := json.NewDecoder(request.Body).Decode(&order)
	if err != nil {
		http.Error(response, "Unknown JSON arguments", http.StatusUnprocessableEntity)
		return 
	}
	
	err = th.appLogic.CreateOrder(order)
	if err != nil {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(response).Encode(http.StatusOK)
	return 
}


func (th TransactionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	//better use function if contains
	if(!strings.Contains(request.Header.Get("Content-Type"), AcceptedContentType)){
		http.Error(response, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}
	
	var err error
	defer func(){
		if err != nil {
			http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	if request.Method == http.MethodPut {
		th.putMethod(response, request)
	}else if request.Method == http.MethodGet{
		return
		//err = th.getMethod(response, request)
	}else{
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

