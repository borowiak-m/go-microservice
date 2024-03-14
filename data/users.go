package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/borowiak-m/go-microservice/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           primitive.ObjectID `json:"_id"`
	FirstName    string             `json:"first_name" validate:"required"` //, min=2, max=100
	LastName     string             `json:"last_name" validate:"required"`  //, min=2, max=100
	Password     string             `json:"password" validate:"required"`   //, min=6, max=30
	Email        string             `json:"email" validate:"required"`      //, email
	Phone        string             `json:"phone" validate:"required"`
	Token        string             `json:"token"`
	UserType     string             `json:"user_type" validate:"required"` //, eq=ADMIN|eq=USER
	RefreshToken string             `json:"refresh_token"`
	CreatedOn    time.Time          `json:"-"`
	UpdatedOn    time.Time          `json:"-"`
	UserId       string             `json:"user_id"`
}

type UserLogin struct {
	Password string `json:"password" validate:"required"` //, min=6, max=30
	Email    string `json:"email" validate:"required"`    //, email
}

var UserCollection *mongo.Collection = OpenCollection(MongoCfg.Client, MongoCfg.DatabaseName, "user")
var ErrUserNotFound = fmt.Errorf("[Error] User not found")
var ErrUserAlreadyExists = fmt.Errorf("[Error] User already exists")

func GetUserByID(id string) (*User, error) {
	user := User{}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err := UserCollection.FindOne(ctx, bson.M{"user_id": id}).Decode(&user)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	user := User{}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err := UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func UpdateAllUserTokens(token, refreshToken, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateObj primitive.D
	userUpdatedOn := time.Now().UTC()

	updateObj = append(updateObj, bson.E{"token", token})
	updateObj = append(updateObj, bson.E{"refresh_token", refreshToken})
	updateObj = append(updateObj, bson.E{"updated_on", userUpdatedOn})
	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	if err != nil {
		return err
	}
	return nil
}

func AddUser(user *User) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// check if user already exists under this email
	userExists, err := doesUserAlreadyExist(user.Email)
	if err != nil {
		return nil, err
	}
	if userExists {
		return nil, ErrUserAlreadyExists
	}
	user.CreatedOn = time.Now().UTC()
	user.UpdatedOn = time.Now().UTC()
	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()
	token, refreshToken, err := helpers.GenerateAllTokens(
		user.Email,
		user.FirstName,
		user.LastName,
		user.UserType,
		user.UserId,
	)
	if err != nil {
		return nil, errors.New("[Error] generating tokens")
	}
	user.Token = token
	user.RefreshToken = refreshToken
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return nil, errors.New("[Error] hashing user password")
	}
	fmt.Println("[ADD TO DB] user ", user)
	resultInsertionNumber, err := UserCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, errors.New("[Error] inserting user to DB")
	}
	return resultInsertionNumber, nil
}

func hashPassword(pass string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		return string(hashPass), err
	}
	return string(hashPass), nil

}

func doesUserAlreadyExist(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	count, err := UserCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, errors.New("[Error] error checing for email ")
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
