package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
	helper "github.com/borowiak-m/go-microservice/helpers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	log *log.Logger
	val *data.Validation
}

type KeyUser struct{}
type KeyUserLogin struct{}

func NewUsers(log *log.Logger, val *data.Validation) *Users {
	return &Users{log, val}
}

// POST request function to handle a login of a user
func (users *Users) Login(respW http.ResponseWriter, req *http.Request) {
	users.log.Println("POST Login user request response")
	userLogin := req.Context().Value(KeyUserLogin{}).(data.UserLogin)
	users.log.Println("[DEBUG] processing login for user:", userLogin)
	foundUser, err := data.GetUserByEmail(userLogin.Email)
	// handle error
	switch err {
	case nil:
	case data.ErrUserNotFound:
		respW.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	default:
		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	}
	// check if login user password match one in the db
	passwordsMatch, err := verifyUserPassword(userLogin.Password, foundUser.Password)
	if err != nil {
		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	}
	if !passwordsMatch {
		//pass don't match
		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "[Error] password does not match"}, respW)
		return
	}
	token, refreshToken, err := helper.GenerateAllTokens(
		foundUser.Email,
		foundUser.FirstName,
		foundUser.LastName,
		foundUser.UserType,
		foundUser.UserId,
	)
	err = data.UpdateAllUserTokens(token, refreshToken, foundUser.UserId)
	if err != nil {
		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "[Error] failed updating user tokens"}, respW)
		return
	}
	foundUser, err = data.GetUserByID(foundUser.UserId)
	if err != nil {
		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	}
	// success
	respW.WriteHeader(http.StatusOK)
	data.ToJSON(foundUser, respW)
}

// POST request function to handle a creation of a new user
func (users *Users) CreateUser(respW http.ResponseWriter, req *http.Request) {
	users.log.Println("POST Signup user request response")
	user := req.Context().Value(KeyUser{}).(data.User)
	users.log.Println("[DEBUG] inserting user:", user)
	insertNum, err := data.AddUser(&user)
	if err != nil {
		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
	}
	// success
	respW.WriteHeader(http.StatusOK)
	data.ToJSON(insertNum, respW)
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

func getUserId(req *http.Request) string {
	vars := mux.Vars(req)
	return vars["id"]
}

func verifyUserPassword(password1, password2 string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2)); err != nil {
		return false, errors.New("[Error] email or password is incorrect")
	}
	return true, nil
}

// func GetUsers()
