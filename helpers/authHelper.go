package helper

import (
	"errors"
	"log"
	"net/http"
)

const (
	AdminUser    = "ADMIN"
	NonAdminUser = "USER"
)

func CheckUserType(req *http.Request, role string) error {
	userType := getUserType(req)
	if userType != role {
		return errors.New("unauthorised to access this user data")
	}
	return nil
}

func MatchUserTypeToId(req *http.Request, userId int) error {
	userType := getUserType(req)
	uid := getUserUid(req)
	log.Printf("From MatchUserTypeToId, intercepted usertype %s and uid %v", userType, uid)
	if userType == NonAdminUser && uid != userId {
		return errors.New("unauthorised to access this user data")
	}
	return CheckUserType(req, userType)
}

func getUserType(req *http.Request) string {
	//return req.Context().Value("user_type").(string)
	return "ADMIN"
}

func getUserUid(req *http.Request) int {
	//return req.Context().Value("uid").(int)
	return 1234
}
