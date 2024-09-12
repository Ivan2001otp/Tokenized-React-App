package model

import (
	util "Backend/Util"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID 			primitive.ObjectID		`bson:"_id"`
	User_id		string					`json:"user_id"`
	Username 	string 					`json:"user_name"`
	PasswordHash	string				`json:"password_hash"`
	Role			string				`json:"role" validate:"required,eq=USER|eq=ADMIN"`
	Token			*string				`json:"token"`
	Refresh_token	*string				`json:"refresh_token"`
	Created_at		time.Time			`json:"created_at"`		
	Updated_at		time.Time			`json:"updated_at"`
}

type TokenSign struct {
	Refresh_token *string	`json:"refresh_token"`
	Flag			string	`json:"flag"`
}

type TokenClaims struct {
	jwt.StandardClaims
	ID  string  `json:"id"`
	Role string `json:"role"`
	Csrf string `json:"csrf"`

}

const RefreshTokenValidTime = time.Hour*72;
const AuthTokenValidTime = time.Minute*45;

func GenerateCSRFsecret()(string,error){
	return util.GenerateRandomString(32);
}

