package handlers

import (
	"log"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
)

type Users struct {
	log *log.Logger
	val *data.Validation
}

func NewUsers(log *log.Logger, val *data.Validation) *Users {
	return &Users{log, val}
}

func (users *Users) Get200(respW http.ResponseWriter, req *http.Request) {
	respW.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericResponse{Message: "All good from Get200"}, respW)
}

func (users *Users) GetSingleUser(respW http.ResponseWriter, req *http.Request) {
	respW.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericResponse{Message: "All good from GetSingleUser"}, respW)
}

func HashPassword()
func VerifyPassword()
func Signup()
func Login()
func GetUsers()
