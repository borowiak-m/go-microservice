package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/borowiak-m/go-microservice/data"
	"github.com/google/uuid"
)

func CheckUserSessionValid(respW http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			respW.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&GenericError{Message: err.Error()}, respW)
			return
		}
		respW.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	}

	sessionToken := cookie.Value

	// check if session exists
	userSession, exists := sessions[sessionToken]
	if !exists {
		respW.WriteHeader(http.StatusUnauthorized)
		data.ToJSON(&GenericError{Message: fmt.Sprintf("[Error] session does not exist for token %s", sessionToken)}, respW)
		return
	}

	// is session exists but expired = unauthorised access
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		respW.WriteHeader(http.StatusUnauthorized)
		data.ToJSON(&GenericError{Message: fmt.Sprintf("[Error] session expired for token %s", sessionToken)}, respW)
		return
	}
}

func RefreshUserSession(respW http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// no cookie rrror
			return
		}
		// bad request
		return
	}

	sessionToken := cookie.Value

	// check if session exists
	userSession, exists := sessions[sessionToken]
	if !exists {
		// status unauthorised
		return
	}

	// is session exists but expired = unauthorised access
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		// status unauthorised
		return
	}

	newToken := uuid.NewString()
	setTokenExpiry := time.Now().Add(120 * time.Second)

	// add new session
	sessions[newToken] = session{
		userEmail: userSession.userEmail,
		expiry:    setTokenExpiry,
	}

	// delete old session
	delete(sessions, sessionToken)

	// set new session token in a cookie
	http.SetCookie(respW, &http.Cookie{
		Name:    "session_token",
		Value:   newToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func Logout(respW http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// no cookie rrror
			return
		}
		// bad request
		return
	}

	sessionToken := cookie.Value

	// delete old session
	delete(sessions, sessionToken)

	// set session token to empty in a cookie, let client know
	http.SetCookie(respW, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}

func GetUserEmailFromSession(req *http.Request) (string, error) {
	cookie, err := req.Cookie("session_token")
	if err != nil {
		return "", err
	}

	sessionToken := cookie.Value

	// check if session exists
	userSession, exists := sessions[sessionToken]
	if !exists {
		return "", fmt.Errorf("[Error] session does not exist")
	}

	// is session exists but expired = unauthorised access
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		return "", fmt.Errorf("[Error] session expired")
	}

	return userSession.userEmail, nil
}
