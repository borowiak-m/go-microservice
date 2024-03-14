package handlers

import (
	"context"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
)

func (users *Users) MiddlewareUserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respW http.ResponseWriter, req *http.Request) {
		users.log.Println("Hello from users auth middleware")
		// if we need to add values to context to pass the req with
		//ctx := context.WithValue(req.Context(), SomeKey{}, someValue)
		//req = req.WithContext(ctx)

		next.ServeHTTP(respW, req)
	})

}

func (users *Users) MiddlewareUserSignupValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respW http.ResponseWriter, req *http.Request) {
		// define empty user object
		user := data.User{}
		// deserialize json from request body into the empty object
		err := data.FromJSON(&user, req.Body)
		// handler errors
		if err != nil {
			users.log.Println("[ERROR] deserializing User", err)
			respW.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, respW)
			return
		}
		users.log.Println("[DEBUG] processing user:", user)
		// validate the user
		// capture validation errors
		vErrs := users.val.Validate(user)
		if vErrs != nil {
			users.log.Println("[ERROR] validating User", err)
			respW.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: vErrs.Errors()}, respW)
			return
		}
		// add user obj to conext and add to request
		ctx := context.WithValue(req.Context(), KeyUser{}, user)
		req = req.WithContext(ctx)
		// pass to next handler
		next.ServeHTTP(respW, req)
	})
}

func (users *Users) MiddlewareUserLoginValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respW http.ResponseWriter, req *http.Request) {
		// define empty user object
		userLogin := data.UserLogin{}
		// deserialize json from request body into the empty object
		err := data.FromJSON(&userLogin, req.Body)
		// handler errors
		if err != nil {
			users.log.Println("[ERROR] deserializing User Login", err)
			respW.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, respW)
			return
		}
		users.log.Println("[DEBUG] processing user login:", userLogin)
		// validate the user
		// capture validation errors
		vErrs := users.val.Validate(userLogin)
		if vErrs != nil {
			users.log.Println("[ERROR] validating User Login", err)
			respW.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: vErrs.Errors()}, respW)
			return
		}
		// add user obj to conext and add to request
		ctx := context.WithValue(req.Context(), KeyUserLogin{}, userLogin)
		req = req.WithContext(ctx)
		// pass to next handler
		next.ServeHTTP(respW, req)
	})
}
