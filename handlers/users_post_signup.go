package handlers

import (
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
)

// POST request function to handle a creation of a new user
func (users *Users) Signup(respW http.ResponseWriter, req *http.Request) {
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
