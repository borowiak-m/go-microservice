package handlers

import (
	"fmt"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
	helper "github.com/borowiak-m/go-microservice/helpers"
)

func (users *Users) GetSingleUser(respW http.ResponseWriter, req *http.Request) {
	userId := getUserId(req)

	users.log.Println("[DEBUG] get user id:", userId)

	if err := helper.MatchUserTypeToId(req, userId); err != nil {
		respW.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{
			Message: fmt.Sprintf("[ERROR] getting user id:%v with error:%s",
				userId,
				err.Error(),
			)}, respW)
		return
	}

	user, err := data.GetUserByID(userId)
	switch err {
	case nil:
	case data.ErrUserNotFound:
		users.log.Println("[ERROR] fetching user:", err)
		respW.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{
			Message: fmt.Sprintf("[ERROR] getting user id:%v with error:%s",
				userId,
				err.Error(),
			)}, respW)
		return
	default:
		users.log.Println("[ERROR] fetching user:", err)

		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	}

	if err = data.ToJSON(user, respW); err != nil {
		users.log.Println("[ERROR] serializing user", err)
	}
}
