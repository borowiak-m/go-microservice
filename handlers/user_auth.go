package handlers

import "net/http"

func (users *Users) MiddlewareUserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respW http.ResponseWriter, req *http.Request) {
		users.log.Println("Hello from users auth middleware")
	})
}
