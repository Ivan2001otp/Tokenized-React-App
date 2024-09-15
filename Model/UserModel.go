package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID	primitive.ObjectID	`bson:"_id"`
	First_name string 	`json:"first_name" validate:"required,min=3,max=30"`
	Last_name  string	`json:"last_name"  validate:"required,min=3,max=30"`
	User_name 	string	`json:"user_name"  validate:"required,min=3,max=30"`
	Email		string		`json:"email"  validate:"required,email"`
	Password	string		`json:"password" validate:"required,min=3"`
	Phone		string		`json:"phone"	validate:"required"`

	Auth_token		 string 		`json:"auth_token"`
	Refresh_token	 string	`json:"refresh_token"`
	User_type 	string		`json:"user_type" validate:"required,eq=USER|eq=ADMIN"`
	User_ID	    string		`json:"user_id"`
	Created_at	time.Time		`json:"created_at"`
	Updated_at	time.Time		`json:"updated_at"`
}