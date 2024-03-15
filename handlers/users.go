package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const (
	AdminUser    = "ADMIN"
	NonAdminUser = "USER"
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

func getUserIdFromURI(req *http.Request) string {
	vars := mux.Vars(req)
	return vars["user_id"]
}

func verifyUserPassword(password1, password2 string) (bool, error) {
	log.Println("[DEBUG] passed password 1:", password1, ", passed password 2:", password2)
	if err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2)); err != nil {
		return false, err
	}
	return true, nil
}

func GetUserTypeFromLoggedInUser(req *http.Request) (string, error) {
	return getUserDataFieldFromLoggedInUser(req, "user_type")
}

func GetUserUserIdFromLoggedInUser(req *http.Request) (string, error) {
	return getUserDataFieldFromLoggedInUser(req, "user_id")
}

func getUserDataFieldFromLoggedInUser(req *http.Request, fieldName string) (string, error) {
	userEmail, err := GetUserEmailFromSession(req)
	if err != nil {
		return "", err
	}
	user, err := data.GetUserByEmail(userEmail)
	if err != nil {
		return "", err
	}
	switch fieldName {
	case "user_type":
		return user.UserType, nil
	case "user_id":
		return user.UserId, nil
	}
	return "", fmt.Errorf("[Error] user data field unrecognised:", fieldName)
}

func MatchUserTypeToId(req *http.Request, userIdFromGetReq string) error {
	loggedInUserType, err := GetUserTypeFromLoggedInUser(req)
	if err != nil {
		return err
	}
	loggedInUserId, err := GetUserUserIdFromLoggedInUser(req)
	if err != nil {
		return err
	}
	log.Printf("Logged in user type %s and user id %v", loggedInUserType, loggedInUserId)
	if loggedInUserType == NonAdminUser && loggedInUserId != userIdFromGetReq {
		return errors.New("unauthorised to access this user data")
	}
	return CheckUserType(req, loggedInUserType)
}

func CheckUserType(req *http.Request, role string) error {
	log.Println("[DEBUG] Checking user type")
	userType, err := GetUserTypeFromLoggedInUser(req)
	if err != nil {
		return err
	}
	if userType != role {
		return errors.New("unauthorised to access this user data")
	}
	return nil
}
