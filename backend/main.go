package main

import (
	database "Backend/Database"
	"Backend/Middleware"
	"Backend/controller"
	"Backend/helper"
	"Backend/shared"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)



var userCollection *mongo.Collection = database.GetCollectionByName(shared.USERS);
func main(){

	fmt.Println("hi")

	err := helper.InitJWT()
// database.Connect();
defer database.Close();


if err!=nil{
	log.Fatal("Error initializing the JWT:",err)
}

router := chi.NewRouter();

router.Use(middleware.CorsMiddleware);


router.Group(func(r chi.Router){
	r.Post("/register",controller.SignUp())
	r.Post("/login",controller.SignIn())
})

router.Group(func(r chi.Router){
	r.Use(middleware.Authenticator)
	log.Println("authenticator")
	r.Get("/dashboard",controller.Dashboard())
	r.Get("/logout",controller.SignOut())
	r.Get("/deleteUser",controller.DeleteUser())
})

err = http.ListenAndServe(":8000",router)
if err!=nil{
	log.Fatal("Error starting the server : ",err)
}

}