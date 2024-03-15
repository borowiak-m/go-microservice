package routes

import (
	"log"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
	"github.com/borowiak-m/go-microservice/handlers"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(newlogger *log.Logger, newValidation *data.Validation, muxRouter *mux.Router) {
	//   create handler for users with logger
	handlerUsers := handlers.NewUsers(newlogger, newValidation)
	// create subrouters per method
	getUserRouter := muxRouter.Methods(http.MethodGet).Subrouter()
	postUserSignupRouter := muxRouter.Methods(http.MethodPost).Subrouter()
	postUserLoginRouter := muxRouter.Methods(http.MethodPost).Subrouter()
	// GET /users
	getUserRouter.HandleFunc("/users", handlerUsers.GetAllUsers)
	// GET /users/{id}
	getUserRouter.HandleFunc("/users/{user_id:[0-9]+}", handlerUsers.GetSingleUser)
	getUserRouter.Use(handlerUsers.MiddlewareUserAuth)
	// POST /users/signup
	postUserSignupRouter.HandleFunc("/users/signup", handlerUsers.Signup)
	// POST /users/signup Middleware: executes before it can go to the HandleFunc
	postUserSignupRouter.Use(handlerUsers.MiddlewareUserAuth)
	postUserSignupRouter.Use(handlerUsers.MiddlewareUserSignupValidation)
	// POST /users/login
	postUserLoginRouter.HandleFunc("/users/login", handlerUsers.Login)
	// POST /users/login Middleware: executes before it can go to the HandleFunc
	postUserLoginRouter.Use(handlerUsers.MiddlewareUserAuth)
	postUserLoginRouter.Use(handlerUsers.MiddlewareUserLoginValidation)
}
