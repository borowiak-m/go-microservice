package handlers

import (
	"net/http"
	"time"

	"github.com/borowiak-m/go-microservice/data"
	helper "github.com/borowiak-m/go-microservice/helpers"
)

var sessions = map[string]session{}

// each session contains the user email and session expiry
type session struct {
	userEmail string
	expiry    time.Time
}

// determine if session is expired
func (sesh session) isExpired() bool {
	return sesh.expiry.Before(time.Now())
}

// POST request function to handle a login of a user
func (users *Users) Login(respW http.ResponseWriter, req *http.Request) {
	users.log.Println("POST Login user request response")
	userLogin := req.Context().Value(KeyUserLogin{}).(data.Credentials)
	users.log.Println("[DEBUG] processing login for user:", userLogin)
	foundUser, err := data.GetUserByEmail(userLogin.Email)
	users.log.Println("[DEBUG] found user:", foundUser, "with err?", err)
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
	passwordsMatch, err := verifyUserPassword(foundUser.Password, userLogin.Password)
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
	// store session
	sessionToken := foundUser.Token
	expiresAt := time.Now().Add(120 * time.Second)
	sessions[sessionToken] = session{
		userEmail: foundUser.Email,
		expiry:    expiresAt,
	}
	// store session in a cookie
	http.SetCookie(respW, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	// respond
	respW.WriteHeader(http.StatusOK)
	data.ToJSON(foundUser, respW)
}
