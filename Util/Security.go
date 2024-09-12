package util

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)



func GenerateBcryptHash(password string)(string ,error){
	hash,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost);

	if err!=nil{
		log.Println("Something went wrong in GenerateBcryptHash");
		return "",err;
	}

	return string(hash[:]),nil;
}

func CheckPasswordAgainstHash(hash string,password string)error{
	return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
}