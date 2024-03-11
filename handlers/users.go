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
	data.ToJSON(&GenericResponse{Message: "All good"}, respW)
}
