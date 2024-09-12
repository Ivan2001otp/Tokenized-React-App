package controller

import (
	models "Backend/Model"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveUser(collection *mongo.Collection) http.Handler{


	fn := func(w http.ResponseWriter, r *http.Request){
		
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		
		var user models.User;

		err := json.NewDecoder(r.Body).Decode(&user)

		if err!=nil{
			defer cancel();

			log.Println("Something wrong in SaveUser-1");
			w.WriteHeader(http.StatusBadRequest);
			w.Write([]byte("Request Body syntax is wrong !"))
			http.Error(w,"Request Body syntax in invalid!",http.StatusBadRequest)
			return;
		}

		//check if the user already exists
		existUser := models.User{}

		err = collection.FindOne(ctx,bson.M{"user_name":user.Username}).Decode(&existUser);

		if err==nil{
			defer cancel();
			http.Error(w,"User already exists",http.StatusConflict);
			return;
		}

		defer cancel();

		user.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID();
		user.User_id = user.ID.Hex()
		

	}

	return http.HandlerFunc(fn)
}	

