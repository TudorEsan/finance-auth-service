package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username   *string            `json:"username" bson:"username" binding:"required,min=4,max=20"`
	Email      *string            `json:"email" bson:"email" binding:"required,email"`
	Password   *string            `json:"password,omitempty" binding:"required,min=8,max=50"`
	CreateDate time.Time          `json:"createDate" bson:"createDate"`
}

type UserLoginForm struct {
	Username *string `json:"username" binding:"required,min=3,max=20"`
	Password *string `json:"password" binding:"required,min=8,max=20"`
}

func (user User) String() string {
	return fmt.Sprintf("USER:\nID: %s\n username: %s\n email: %s\n", user.ID.Hex(), *user.Username, *user.Email)
}
