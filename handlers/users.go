package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/borowiak-m/go-microservice/data"
	helper "github.com/borowiak-m/go-microservice/helpers"
	"github.com/gorilla/mux"
)

type Users struct {
	log *log.Logger
	val *data.Validation
}

//type KeyUserType struct{}

func NewUsers(log *log.Logger, val *data.Validation) *Users {
	return &Users{log, val}
}

func (users *Users) Get200(respW http.ResponseWriter, req *http.Request) {
	respW.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericResponse{Message: "All good from Get200"}, respW)
}

func (users *Users) GetSingleUser(respW http.ResponseWriter, req *http.Request) {
	userId := getUserId(req)

	users.log.Println("[DEBUG] get user id:", userId)

	if err := helper.MatchUserTypeToId(req, userId); err != nil {
		respW.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{
			Message: fmt.Sprintf("[ERROR] getting user id:%v with error:%s",
				userId,
				err.Error(),
			)}, respW)
		return
	}

	user, err := data.GetUserByID(userId)
	switch err {
	case nil:
	case data.ErrUserNotFound:
		users.log.Println("[ERROR] fetching user:", err)
		respW.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{
			Message: fmt.Sprintf("[ERROR] getting user id:%v with error:%s",
				userId,
				err.Error(),
			)}, respW)
		return
	default:
		users.log.Println("[ERROR] fetching user:", err)

		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	}

	if err = data.ToJSON(user, respW); err != nil {
		users.log.Println("[ERROR] serializing user", err)
	}

}

func getUserId(req *http.Request) int {
	// parse id from url
	vars := mux.Vars(req)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	return id
}

// func HashPassword()
// func VerifyPassword()
// func Signup()
// func Login()
// func GetUsers()
