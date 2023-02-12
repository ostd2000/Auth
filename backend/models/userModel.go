package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id"`
    UserID string `json:"user_id"`
	
	Username *string `json:"username" validate:"required, unique"`
    PhoneNumber *string `json:"phone_number" validate:"required"`  	

	// Email is not required.
    Email *string  `json:"email" validate:"email"` 

	Password *string `json:"password" validate:"required, min=8"`

	AccessToken *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`

    CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
