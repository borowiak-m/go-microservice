package data

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID           primitive.ObjectID `json:"_id"`
	FirstName    string             `json:"first_name" validate:"required, min=2, max=100"`
	LastName     string             `json:"last_name" validate:"required, min=2, max=100"`
	Password     string             `json:"password" validate:"required, min=6, max=30"`
	Email        string             `json:"email" validate:"required, email"`
	Phone        string             `json:"phone" validate:"required"`
	Token        string             `json:"token"`
	UserType     string             `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	RefreshToken string             `json:"refresh_token"`
	CreatedOn    time.Time          `json:"-"`
	UpdatedOn    time.Time          `json:"-"`
	UserId       string             `json:"user_id"`
}

var UserCollection *mongo.Collection = OpenCollection(MongoCfg.Client, MongoCfg.DatabaseName, "user")
var ErrUserNotFound = fmt.Errorf("User not found")

func GetUserByID(id int) (*User, error) {
	user := User{}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err := UserCollection.FindOne(ctx, bson.M{"user_id": id}).Decode(&user)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}
