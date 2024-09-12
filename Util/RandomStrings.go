package util

import (
	"encoding/base64"
	"log"
	"math/rand"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	if err!=nil{
		log.Println("GenerateRandomBytes->error in reading !")
		return nil,err;
	}

	return b,nil;
}

func GenerateRandomString(s int)(string,error){
	b,err := GenerateRandomBytes(s);
	if err!=nil{
		log.Println("GenerateRandomString->something went wrong");
		return "",err;
	}

	return base64.URLEncoding.EncodeToString(b),nil;
}