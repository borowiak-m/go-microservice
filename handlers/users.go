package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
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

func (users *Users) Get200(respW http.ResponseWriter, req *http.Request) {
	respW.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericResponse{Message: "All good from Get200"}, respW)
}

func getUserId(req *http.Request) string {
	vars := mux.Vars(req)
	return vars["id"]
}

func verifyUserPassword(password1, password2 string) (bool, error) {
	log.Println("[DEBUG] passed password 1:", password1, ", passed password 2:", password2)
	if err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2)); err != nil {
		return false, errors.New("[Error] email or password is incorrect")
	}
	return true, nil
}

// func GetUsers()
