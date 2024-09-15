package main

import (
	"Backend/Middleware"
	"Backend/controller"
	"Backend/helper"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main(){
fmt.Println("hi")

err := helper.InitJWT()


if err!=nil{
	log.Fatal("Error initializing the JWT:",err)
}

router := chi.NewRouter();

router.Group(func(r chi.Router){
	r.Post("/register",controller.SignUp())
	r.Post("/login",controller.SignIn())
})

router.Group(func(r chi.Router){
	r.Use(middleware.Authenticator)

	r.Get("/dashboard",controller.Dashboard())
	r.Get("/logout",controller.SignOut())
	r.Get("/deleteUser",controller.DeleteUser())
})

err = http.ListenAndServe(":8000",router)
if err!=nil{
	log.Fatal("Error starting the server : ",err)
}
}