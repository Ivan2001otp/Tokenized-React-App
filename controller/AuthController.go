package controller

import (
	database "Backend/Database"
	"Backend/Model"
	"Backend/helper"
	"Backend/shared"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type status map[string]interface{};
var validate = validator.New();

func Dashboard()http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		w.Write([]byte("Hello its working!"));
	}
}

func SignUp() http.HandlerFunc{
	return func (w http.ResponseWriter,r *http.Request)  {
		var user model.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err!=nil{
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error":"Invalid json format . "+err.Error()})
			return;
		}

		err = validate.Struct(user)

		if err!=nil{
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error":"Failed to validate json . "+err.Error()})
			return;
		}

		var user_or_email_Count int64;

		user_or_email_Count,err = database.FetchUserCountByCredential(shared.USERS,user.Email,user.User_name)
		
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError);
			json.NewEncoder(w).Encode(status{"error":"Failed to count user-name and email"})
			return;
		}

		if user_or_email_Count>0{
			w.WriteHeader(http.StatusInternalServerError);
			json.NewEncoder(w).Encode(status{"error":"duplicate user_name or email"})
			return;
		}

		if user.User_name==""{
			user.User_name = user.First_name+"-"+user.Last_name
		}

		user.ID = primitive.NewObjectID()
		user.Password = helper.HashPassword(user.Password)
		user.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.User_ID = user.ID.Hex();

		authTokenString,refreshTokenString,csrfString,err := helper.CreateNewTokens(user)
		
		if err!=nil{
			log.Println("SignUp->err1")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error":err.Error()})
			return;
		}

		user.Auth_token = authTokenString
		user.Refresh_token = refreshTokenString
		
		//save user entity...
		result,err := database.SaveUserCredential(shared.USERS,user)

		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error":"Failed to insert data !"})
			return;
		}

		log.Println("Saved userid ",result)

		helper.SetAuthAndRefreshCookies(&w,authTokenString,refreshTokenString);
		w.Header().Set("X-CSRF-Token",csrfString)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status{"InsertID":user.User_ID})

	}
}


//login
func SignIn() http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		//take only email and password during login.
		var user model.User;

		if err:= json.NewDecoder(r.Body).Decode(&user);err!=nil{
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error":"invalid json formt . "+err.Error()})
			return;
		}

		if user.Email=="" || user.Password==""{
			w.WriteHeader(http.StatusBadRequest);
			json.NewEncoder(w).Encode(status{"error":"Enter email or password field"})
			return;
		}

		var foundUser *model.User
		//fetch the user by mail.
		foundUser,err1 := database.FetchUserBySpecificCredential(shared.USERS,user.Email)

		if err1!=nil{
			w.WriteHeader(http.StatusInternalServerError);
			json.NewEncoder(w).Encode(status{"error":err1.Error()})
			return;
		}

		err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password),[]byte(user.Password))

		if err!=nil{
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(status{"error":err.Error()})
				return;
		}

		authToken,refToken,csrfString,err := helper.CreateNewTokens(*foundUser);

		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error":err.Error()});
			return;
		}

		helper.SetAuthAndRefreshCookies(&w,authToken,refToken);

		w.Header().Set("X-CSRF-Token",csrfString);

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(foundUser)
	}
}


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
		return "",errors.New("error occured while fetching authToken");
	})

	if err!=nil{
		return "","","",errors.New("error occured while fetching authToken")
	}

	authTokenClaims,ok := authToken.Claims.(*helper.SignedDetails)

	if !ok{
		return "","","",errors.New("grabUser->Error while fetching authToken claims")

	}

	return authTokenClaims.Last_name,authTokenClaims.User_name,authTokenClaims.Email,nil;
}