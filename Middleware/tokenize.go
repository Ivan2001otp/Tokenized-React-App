package middleware

import (
	models "Backend/Model"
	util "Backend/Util"
	db "Backend/Database"
	shared "Backend/shared"
	"context"
	"crypto/rsa"
	"errors"
	"log"
	"time"
	"io/ioutil"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	privateKeyPath = "key/app.rsa"
	publicKeyPath = "key/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey *rsa.PrivateKey
)

func InitJwt() error{
	signBytes,err := ioutil.ReadFile(privateKeyPath);

	if err!=nil{
		log.Println("init Jwt went wrong 1");
		return err;
	}

	signKey,err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err!=nil{
		log.Println("init Jwt went wrong 2");
		return err;
	}

	log.Println("SignKey generated")
	log.Println(signKey);

	verifyBytes,err := ioutil.ReadFile(publicKeyPath);
	if err!=nil{
		log.Println("init Jwt went wrong 3");
		return err;
	}

	verifyKey,err := jwt.ParseECPublicKeyFromPEM(verifyBytes)

	if err!=nil{
		log.Println("init Jwt went wrong 4")
		return err;
	}

	log.Println("Verify key generated!");
	log.Println(verifyKey)
	return nil;
}

func RevokeRefreshToken(collection *mongo.Collection,refreshTokenString string) error{
	refreshToken,err := jwt.ParseWithClaims(refreshTokenString,&models.TokenClaims{},func(token *jwt.Token)(interface{},error){
		return verifyKey,nil;
	})

	if err!=nil{
		return errors.New("Could not parse token with claims !")
	}

	refreshTokenClaims,ok := refreshToken.Claims.(*models.TokenClaims)

	if !ok{
		return errors.New("Could not fetch refreshTokenClaims !")
	}

	//delete refreshtoken
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

	filter := bson.M{"refresh_token":refreshTokenClaims.Id}
	update := bson.M{"$set":bson.M{"refresh_token":""}}
	updateResult,err := collection.UpdateOne(ctx,
		filter,update,
	)

	if err!=nil{
		defer cancel();
		return errors.New("Failed to revoke refresh token")
	}

	defer cancel();
	//check if update was successful
	if updateResult.MatchedCount==0{
		
		return errors.New("Refresh Token not found!");
	}

	log.Println("Revoked refresh token successfully !");
	return nil;
}

func CreateNewTokens(uuid string,role string)(authTokenString string,refreshTokenString string,csrfSecret string,err error){
	csrfSecret,err = models.GenerateCSRFsecret()

	if err!=nil{
		log.Println("CreateNewTokens threw error while generating csrf");
		return "","","",err;
	}

	refreshTokenString,err = CreateRefreshTokenString(uuid,role,csrfSecret);

	if err!=nil{
		log.Println("CreateRefreshTokenString threw error . Could not generate refrshToken!");
		return "","","",err;
	}

	authTokenString,err = 	CreateAuthTokenString(uuid,role,csrfSecret);

	if err!=nil{
		log.Println("CreateAuthTokenString threw error . Could not generate authToken!");
		return "","","",err;
	}

	return authTokenString,refreshTokenString,csrfSecret,nil;
}


//create the auth token
func CreateAuthTokenString(uuid string,role string,csrfKey string)(authTokenString string,err error){
	authTokenExpiry := time.Now().Add(models.AuthTokenValidTime).Unix()

	authClaims := models.TokenClaims{
		jwt.StandardClaims{
			Subject: uuid,
			ExpiresAt: authTokenExpiry,
		},
		uuid,
		role,
		csrfKey,
	}

	authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"),authClaims);

	authTokenString,err = authJwt.SignedString(signKey);
	if err!=nil{
		log.Println("CreateAuthTokenString threw error . authJwt went wrong!")
		return "",err;
	}

	return authTokenString,nil;
}


//create the refresh token
func CreateRefreshTokenString(uuid string,role string,csrfKey string)(refreshTokenString string,err error){
	refreshTokenExpire := time.Now().Add(models.RefreshTokenValidTime).Unix();

	refreshJti,err := util.GenerateRandomString(32);


	if err!=nil{
		log.Println("CreateRefreshTokenString->refreshTokenJTi could not be generated!")
		return "",err;
	}

	err =  db.SaveRefreshToken(shared.RefreshTokenIDs,refreshJti);

	if err!=nil{
		log.Println("CreateRefreshTokenString-something went wrong while saving refresh jti token !")
		return "",err;
	}

	refreshClaims := models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id: refreshJti,
			Subject: uuid,
			ExpiresAt: refreshTokenExpire,
		},
		ID: refreshJti,
		Role: role,
		Csrf: csrfKey,
	}

	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"),refreshClaims)

	newRefreshTokenString,err := refreshJwt.SignedString(signKey);

	if err!=nil{
		log.Println("CreateRefreshTokenString -> new refresh token string could not be created!");
		return "",err;
	}

	return newRefreshTokenString,nil;
}

/*
To get the uuid ,first get the authtoken.
through token get the authclaims.
From authclaims get the uuid which is present in subject.
*/
func GrabUUID(authTokenString string)(string,error){
	authToken,_ := jwt.ParseWithClaims(authTokenString,&models.TokenClaims{},func(token *jwt.Token)(interface{},error){
		return "",errors.New("Error fetching claims -1")
	})

	authClaims,ok := authToken.Claims.(*models.TokenClaims)

	if !ok{
		return "",errors.New("Error fetching claims -2")
	}

	return authClaims.StandardClaims.Subject,nil;
}

//update refresh token
func UpdateRefreshToken(oldRefreshTokenString string)(newRefreshTokenString string,err error){
	refreshToken,err := jwt.ParseWithClaims(oldRefreshTokenString,&models.TokenClaims{},func(token *jwt.Token)(interface{},error){
		return verifyKey,nil;
	})

	refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()

	oldRefreshTokenClaims,ok := refreshToken.Claims.(*models.TokenClaims)

	if !ok{
		return "",errors.New("UpdateRefreshToken->refresh claims threw error")
	}

	refreshClaims := models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id: oldRefreshTokenClaims.Id,
			ExpiresAt: refreshTokenExp,
		},
		ID: oldRefreshTokenClaims.ID,
		Role: oldRefreshTokenClaims.Role,
		Csrf: oldRefreshTokenClaims.Csrf,
	}

	refreshJwt:=jwt.NewWithClaims(jwt.GetSigningMethod("RS256"),refreshClaims)

	newRefreshTokenString,err = refreshJwt.SignedString(refreshJwt)
	
	if err!=nil{
		log.Println("UpdateRefreshToken-> new refresh token could not get generated!")
		return "",err;
	}

	return newRefreshTokenString,nil;
}

//update authentication token
/*
First we will check whether refresh token is active or not.
if the refresh token is active only at that time we will update the token.
In return we get new refresh token ,new csrf token
*/
func UpdateAuthToken(refreshTokenString string,oldAuthTokenString string)(newAuthTokenString string,csrfSecret string,err error){
	
	refreshToken,err := jwt.ParseWithClaims(refreshTokenString,
		&models.TokenClaims{},
		func(token *jwt.Token)(interface{},error){
		return verifyKey,nil
	})

	refreshClaims,ok := refreshToken.Claims.(*models.TokenClaims)

	if !ok{
		err = errors.New("UpdateAuthToken->Error reading refreshClaims");
		return "","",err;
	}

	 tokenExists,err := db.CheckRefreshTokenAlreadyexists(shared.RefreshTokenIDs,refreshTokenString);

	 if err!=nil{
		log.Println("UpdateAuthToken->Error reading refreshClaims 1")
		return "","",err;
	 }else{
		if (tokenExists){
			//if refresh token exists
			log.Println("refrsh token exists!");
		}else{
		log.Println("UpdateAuthToken->Error reading refreshClaims 2")

			//if it does not exists..
			return "","",errors.New("Unauthorized")
		}
	 }

	 if refreshToken.Valid{
		//if the refresh token has not expired yet
		authToken,_ := jwt.ParseWithClaims(oldAuthTokenString,&models.TokenClaims{},func(token *jwt.Token)(interface{},error){
			return verifyKey,nil;
		})

		oldAuthTokenClaims,ok:=authToken.Claims.(*models.TokenClaims)

		if !ok{
		log.Println("UpdateAuthToken->Error reading auth claims 3")
			err = errors.New("Error reading auth claims")
			return "","",err;
		}

		newcsrfSecret,err := models.GenerateCSRFsecret();

		if err!=nil{
		log.Println("UpdateAuthToken->Error in creating csrf secret 4")
			
			return "","",err;
		}
	 
	
		newAuthTokenString,err = CreateAuthTokenString(oldAuthTokenClaims.Subject,oldAuthTokenClaims.Role,oldAuthTokenClaims.Csrf);

		if err!=nil{
		log.Println("UpdateAuthToken->Error 5")

			log.Panic(err)
			return "","",err;
		}

		return newAuthTokenString,newcsrfSecret,nil;
	}

	//if refresh token is expired ,delete it.
	log.Println("Refresh token expired")
	db.DeleteRefreshToken(shared.RefreshTokenIDs,refreshClaims.StandardClaims.Id);
	return "","",errors.New("Unauthorized!")
}