package handlers

import (
	"net/http"
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
