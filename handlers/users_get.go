package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/borowiak-m/go-microservice/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (users *Users) GetSingleUser(respW http.ResponseWriter, req *http.Request) {
	userId := getUserIdFromURI(req)

	users.log.Println("[DEBUG] get user id:", userId)

	if err := MatchUserTypeToId(req, userId); err != nil {
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

func (users *Users) GetAllUsers(respW http.ResponseWriter, req *http.Request) {
	// check if user is ADMIN, only ADMIN can get all users
	if err := CheckUserType(req, AdminUser); err != nil {
		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: fmt.Sprintf("[Error] getting user type; %s", err)}, respW)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// [TO DO] Server-side pagination
	// get records per page definition
	// get page index (or default to start index) from request
	startIndex := 1
	recordPerPage := 1000
	// db search empty filter
	// [TO DO] abstract into a user search function
	matchStage := bson.D{{"$match", bson.D{{}}}}
	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{{"_id", "null"}}},
		{"total_count", bson.D{{"$sum", 1}}},
		{"data", bson.D{{"$push", "$$ROOT"}}}}}}
	projecStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}
	result, err := data.UserCollection.Aggregate(ctx, mongo.Pipeline{groupStage, matchStage, projecStage})
	if err != nil {
		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: fmt.Sprintf("[Error] getting user type; %s", err)}, respW)
	}

	var allUsers []bson.M
	if err = result.All(ctx, &allUsers); err != nil {
		log.Fatal(err)
	}

	respW.WriteHeader(http.StatusOK)
	data.ToJSON(allUsers[0], respW)
}
