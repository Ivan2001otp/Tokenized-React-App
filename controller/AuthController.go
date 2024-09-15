package controller

import (
	database "Backend/Database"
	"Backend/helper"
	"Backend/shared"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
)

type status map[string]interface{};
var validate = validator.New();

func SignOut() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("deleted cookies")
	 	helper.NullifyTokenCookies(&w,r);
	}
}

func DeleteUser() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Deleting user")
		authCookie,authErr := r.Cookie(shared.AUTH_TOKEN)

		if authErr == http.ErrNoCookie{
			log.Println("Unauthorized attempt!No auth cookie");
			helper.NullifyTokenCookies(&w,r)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(status{"error":"No auth Cookie. Unauthorized attempt!"})
			return;
		}else if authErr!=nil{
			log.Panicf("panic:+%v",authErr)
			helper.NullifyTokenCookies(&w,r)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error":authErr.Error()})
		}

		last_name,user_name,email,err := grabUser(authCookie.Value)

		if err!=nil{
			log.Panic(err)
			helper.NullifyTokenCookies(&w,r);
			w.WriteHeader(http.StatusInternalServerError);
			log.Println("DeleteUser->grab user threw error!")
			json.NewEncoder(w).Encode(status{"error":err.Error()})
			return;
		}

		//delete the user record from mongo
		err = database.DeleteUserByCredentials(shared.USERS,user_name,last_name,email)

		if err!=nil{
			helper.NullifyTokenCookies(&w,r);
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error":"Database error while deletion : "+err.Error()})
			return;
		}

		helper.NullifyTokenCookies(&w,r);
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status{"msg":"Deleted user "+user_name});
	}
}

func grabUser(authTokenString string)(last_name,user_name,email string,err error){
	authToken,err := jwt.ParseWithClaims(authTokenString,&helper.SignedDetails{},func(token *jwt.Token)(interface{},error){
		return "",errors.New("Error occured while fetching AuthToken!");
	})

	if err!=nil{
		return "","","",errors.New("Error occured while fetching AuthToken!")
	}

	authTokenClaims,ok := authToken.Claims.(*helper.SignedDetails)

	if !ok{
		return "","","",errors.New("grabUser->Error while fetching authToken claims")

	}

	return authTokenClaims.Last_name,authTokenClaims.User_name,authTokenClaims.Email,nil;
}