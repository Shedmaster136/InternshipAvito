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
		err = wh.putMethod(response, request)
	}else if request.Method == http.MethodGet{
		err = wh.getMethod(response, request)
	}else{
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (wh WalletHandler) putMethod(response http.ResponseWriter,request *http.Request) error{
	
	user := businesslogic.User{}
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		//log
		return err
	}
	
	err = wh.appLogic.TopUpWallet(user)
	if err != nil {
		//log
		return err
	}
	
	return nil
}



func (wh WalletHandler) getMethod(response http.ResponseWriter,request *http.Request) error{
	
	fmt.Fprint(response, "Called Get handler\n")
	var userID int64
	err := json.NewDecoder(request.Body).Decode(&userID)
	if err != nil {
		//log
		return err
	}
	
	balance, err := wh.appLogic.GetUserBalance(userID)
	if err != nil {
		//log
		return err
	}
	err = json.NewEncoder(response).Encode(balance)
	return nil
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
