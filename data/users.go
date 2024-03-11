package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
