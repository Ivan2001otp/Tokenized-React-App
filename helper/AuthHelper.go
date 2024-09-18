package helper

import (
	"Backend/Database"
	"Backend/Model"
	"Backend/shared"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTokenValidTime = time.Hour * 72
	authTokenValidTime    = time.Minute * 15
	privateKeyPath        = "keys/private_key.pem"
	publicKeyPath         = "keys/public_key.pem"
	emptyString           = ""
)

var (
	signKey *rsa.PrivateKey
	verifyKey *rsa.PublicKey
)

func InitJWT() error{
	signBytes,err := os.ReadFile(privateKeyPath)

	if err!=nil{
		log.Println("initJWT-1 throwed error")
		log.Println(err)

		return err;
	}

	signKey,err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err!=nil{
		log.Println("Error while parsing the signKey!");
		log.Println(err)
		return err;
	}

	verifyBytes,err := os.ReadFile(publicKeyPath)

	if err!=nil{
		log.Println("Error while reading public key path!")
		log.Println(err)
		return err;
	}

	verifyKey,err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	if err!=nil{
		log.Println("Error while parsing public key from PEM!");
		log.Println(err);
		return err;
	}

	log.Println("VerifyKey generated ",)
	log.Println("Signkey generated",);

	return nil;
}

type SignedDetails struct{
	First_name string
	Last_name string
	User_name string
	Email	  string
	User_ID   string
	User_type 	string
	CSRFtoken   string
	jwt.RegisteredClaims
}

func updateRefreshTokenExpire(oldRefreshTokenString string)(newRefreshTokenString string,err error){
	refreshToken,err := jwt.ParseWithClaims(oldRefreshTokenString,&SignedDetails{},func(t *jwt.Token)(interface{},error){
		return verifyKey,nil
	})

	if err!=nil{
		return;
	}

	oldRefreshTokenClaims,ok := refreshToken.Claims.(*SignedDetails)

	if !ok{
		log.Println("updateRefreshTokenExpire->err1")
		err = errors.New("Refresh Token is invalid!")
		return;
	}

	refreshClaims := &SignedDetails{
		First_name: oldRefreshTokenClaims.First_name,
		Last_name: oldRefreshTokenClaims.Last_name,
		User_name: oldRefreshTokenClaims.User_name,
		Email: oldRefreshTokenClaims.Email,
		User_ID: oldRefreshTokenClaims.User_ID,
		User_type: oldRefreshTokenClaims.User_type,
		CSRFtoken: oldRefreshTokenClaims.CSRFtoken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenValidTime)),
		},
	}

	newRefreshTokenString,err = jwt.NewWithClaims(jwt.SigningMethodRS256,refreshClaims).SignedString(signKey)
	
	if err!=nil{
		log.Println("updateRefreshTokenExpire->err2")
		log.Fatal(err);
		return "",err;
	}	
	
	return newRefreshTokenString,nil;
}

func updateAuthTokenString(refreshTokenString,oldAuthTokenString string,)(newAuthTokenString,csrfString string,err error){
	refreshToken,err := jwt.ParseWithClaims(refreshTokenString,&SignedDetails{},func(token *jwt.Token)(interface{},error){
		return verifyKey,nil;
	})

	if err!=nil{
		log.Println("updateAuthTokenString->err1")
		log.Fatal(err);
		return;
	}

	refreshTokenClaims,ok := refreshToken.Claims.(*SignedDetails)

	if !ok{
		log.Println("updateAuthTokenString->err2")
		err =errors.New("Error reading jwt claims")
		return;
	}

	var refresh_token string

	refresh_token,err = database.FetchRefreshTokenForAuthtokenUpdate(shared.USERS,
		refreshTokenClaims.User_name,
		refreshTokenClaims.Email,refreshTokenClaims.Last_name);

	if err!=nil{
		log.Println("updateAuthTokenString->err3")
		log.Fatal(err);
		return;
	}

	if refresh_token!=""{
		if refreshToken.Valid{
			authToken,err1 := jwt.ParseWithClaims(oldAuthTokenString,&SignedDetails{},func(t *jwt.Token)(interface{},error){
				return verifyKey,nil
			})

			if err1!=nil{
				log.Println("updateAuthTokenString->err4")
				log.Fatal(err)
				return "","",err;
			}

			oldAuthTokenClaims,ok := authToken.Claims.(*SignedDetails)

			if !ok{
				log.Println("updateAuthTokenString->err5")

				err= errors.New("Error reading jwt claims ")
				return "","",err;
			}
			csrfString,err = generateCSRFSecret();

			if err!=nil{
				log.Println("updateAuthTokenString->err6")

				log.Fatal(err);
			}

			newAuthTokenString,err = createAuthTokenString(oldAuthTokenClaims.First_name,oldAuthTokenClaims.Last_name,
				oldAuthTokenClaims.User_name,oldAuthTokenClaims.Email,
				oldAuthTokenClaims.User_ID,oldAuthTokenClaims.User_type,
				csrfString)
			return newAuthTokenString,csrfString,err;

		}else{
			log.Println("Refresh Token has expired err5")
			log.Println("updateAuthTokenString->err7")
			

			//if refrestoken is valid , then set the refreshtoken as empty
			err3 := database.SetRefreshTokenAsEmpty(shared.USERS,refreshTokenClaims.Last_name,
				refreshTokenClaims.User_name,refreshTokenClaims.Email)
			
			if err3!=nil{

				log.Println("updateAuthTokenString->err8")
				log.Fatal(err3)
				return "","",errors.New("Unauthorized");
			}
			err = errors.New("Unauthorized")
			return;
		}
	}else{
		log.Println("Refresh Token has been revoked!");
		err = errors.New("Unauthorized")
		return;
	}
}

func updateRefreshTokenCSRF(oldRefreshTokenString,newCSRFString string)(newRefreshTokenString string,err error){
	refreshToken,err := jwt.ParseWithClaims(oldRefreshTokenString,&SignedDetails{},func(token *jwt.Token)(interface{},error){
		return verifyKey,nil;
	})

	if err!=nil{
		log.Println("updateRefreshTokenCSRF->err 1")
		log.Fatal(err)
		return;
	}

	oldRefreshTokenClaims,ok := refreshToken.Claims.(*SignedDetails)
	
	if !ok{
		log.Fatalf("updateRefreshTokenCSRF->err 2")
		err = errors.New("Could not read refresh token claims");
		return;
	}

	refreshClaims := &SignedDetails{
		First_name: oldRefreshTokenClaims.First_name,
		Last_name: oldRefreshTokenClaims.Last_name,
		User_name: oldRefreshTokenClaims.User_name,
		Email: oldRefreshTokenClaims.Email,
		User_ID:oldRefreshTokenClaims.User_ID,
		User_type: oldRefreshTokenClaims.User_type,
		CSRFtoken: newCSRFString,
	RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenValidTime)),
	},
	}

	newRefreshTokenString,err = jwt.NewWithClaims(jwt.SigningMethodRS256,refreshClaims).SignedString(signKey)

	if err!=nil{
		log.Println("updateRefreshTokenCSRF->err 3")
		log.Fatal(err)
		return;
	}

	return newRefreshTokenString,nil;
}

func GrabCSRFfromRequest(r *http.Request) string{
	csrfFromFrom := r.FormValue(shared.X_CSRF_Token)

	log.Println(r)
	log.Println("postman-token:",r.Header.Get(shared.POSTMAN_TOKEN))
	
	if csrfFromFrom!=""{
		log.Println("The csrf token is1 ->",csrfFromFrom)
		return csrfFromFrom;
	}else{
		log.Println("The csrf token is2 ->",r.Header.Get(shared.X_CSRF_Token))
		return r.Header.Get(shared.X_CSRF_Token)
	}
}

func HashPassword(password string)string{
	bytes,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)

	if err!=nil{
		panic(err)
	}
	return string(bytes);
}

func SetAuthAndRefreshCookies(w *http.ResponseWriter,authTokenString ,refreshTokenString string){

	authCookie := http.Cookie{
		Name: shared.AUTH_TOKEN,
		Value:authTokenString,
		// SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		// Secure: true,
		Path:"/",
		Domain: "localhost:8000",
	}


	http.SetCookie(*w,&authCookie);


	refreshCookie := http.Cookie{
		Name:shared.REFRESH_TOKEN,
		Value: refreshTokenString,
		HttpOnly: true,
		Path: "/",
		Domain: "localhost:8000",
		// SameSite: http.SameSiteLaxMode,
		// Secure: true,
	}

	http.SetCookie(*w,&refreshCookie);
}

func NullifyTokenCookies(w *http.ResponseWriter,r *http.Request){
	AuthCookie := http.Cookie{
		Name: shared.AUTH_TOKEN,
		Value: "",
		Expires: time.Now().Add(-100*time.Hour),
		HttpOnly: true,
		
	}

	http.SetCookie(*w,&AuthCookie)


	RefreshCookie := http.Cookie{
		Name: shared.REFRESH_TOKEN,
		Value: "",
		Expires: time.Now().Add(-100*time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(*w,&RefreshCookie)
}

func CreateNewTokens(user model.User)(authTokenString,refreshTokenString,csrfString string,err error){
	csrfString,err = generateCSRFSecret()

	if err!=nil{
		log.Println("CreateNewTokens->Could not generate csrf token!")
		log.Panic(err)
		return;
	}

	refreshTokenString,err = createRefreshTokenString(user.First_name,user.Last_name,user.User_name,user.Email,user.User_ID,user.User_type,csrfString);

	if err!=nil{
		log.Println("CreateNewTokens->Could not create refresh token!");
		log.Panic(err);
		return;
	}

	authTokenString,err = createAuthTokenString(user.First_name,user.Last_name,user.User_name,user.Email,user.User_ID,user.User_type,csrfString)

	if err!=nil{
		log.Println("CreateNewTokens->could not create authtoken!");
		log.Panic(err);
		return;
	}

	return authTokenString,refreshTokenString,csrfString,nil;
}

func createRefreshTokenString(first_name,last_name,user_name,email,user_id,
	user_type,csrfToken string)(refreshTokenString string,err error){

		refreshClaims := &SignedDetails{
			First_name: first_name,
			Last_name: last_name,
			User_name: user_name,
			Email: email,
			User_ID: user_id,
			User_type: user_type,
			CSRFtoken: csrfToken,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenValidTime)),
			},
		}

		refreshTokenString,err = jwt.NewWithClaims(jwt.SigningMethodRS256,refreshClaims).SignedString(signKey)

		if err!=nil{
			log.Println("createRefreshTokenString->could not create refresh token!");
			log.Fatal(err)
			return;
		}


		return refreshTokenString,nil;
}

func createAuthTokenString(first_name,last_name,user_name,email,user_id,user_type,
	csrfToken string)(authTokenString string,err error){

		authClaims := &SignedDetails{
			First_name: first_name,
			Last_name: last_name,
			User_name: user_name,
			Email: email,
			User_ID: user_id,
			User_type: user_type,
			CSRFtoken: csrfToken,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(authTokenValidTime)),
			},
		}

		authTokenString,err = jwt.NewWithClaims(jwt.SigningMethodRS256,authClaims).SignedString(signKey);

		if err!=nil{
			log.Println("createAuthTokenString->could not create auth token!")
			log.Fatal(err);
			return;
		}
		return authTokenString,nil;
}	

func generateCSRFSecret()(string,error){
	bytes := make([]byte,32)

	_,err := rand.Read(bytes)

	if err!=nil{
		log.Println("generateCSRFSecret->could not read bytes!")
		log.Println(err);
		return "",err;
	}

	return base64.URLEncoding.EncodeToString(bytes),err;
}

func CheckAndRefreshToken(oldAuthTokenString,
	oldRefreshTokenString,
	oldCSRFString string)(newAuthTokenString,
		newRefreshTokenString,
		newCSRFString string,err error){

			if oldCSRFString==""{
				log.Println("CheckAndRefreshToken->No CSRF token")
				err = errors.New("Unauthorized");
				return;
			}

			authToken,err := jwt.ParseWithClaims(oldAuthTokenString,&SignedDetails{},
				func (token *jwt.Token)(interface{},error){
					return verifyKey,nil;
				});
			
				if err!=nil{
					log.Println("CheckAndRefreshToken->claims jwt could not get")
					return;
				}

				authTokenClaims,ok := authToken.Claims.(*SignedDetails)

				if !ok{
				log.Println("CheckAndRefreshToken->Auth Token is Invalid!")

					err = errors.New("auth token is invalid");
					return "","","",err;
				}

				if oldCSRFString!=authTokenClaims.CSRFtoken{
					log.Println("CheckAndRefreshToken->csrf token do not match jwt");
					err = errors.New("Unauthorized:CSRF token doesn't match jwt");
					return "","","",err;
				}

				if authToken.Valid{
					log.Println("Auth Token is valid")
					newCSRFString = authTokenClaims.CSRFtoken;

					newRefreshTokenString,err = updateRefreshTokenExpire(oldRefreshTokenString)
					newAuthTokenString = oldAuthTokenString
					return
				}else if v,ok := err.(jwt.ClaimsValidator);ok{
					log.Println("CheckAndRefreshToken->Auth Token is not valid")
					expires_at,expire_err := v.GetExpirationTime()

					if expire_err!=nil{
					log.Println("CheckAndRefreshToken->auth token expired")

						panic(expire_err)
						
					}

					if expires_at.Unix() < time.Now().Unix(){
						log.Println("Auth token expired")
						newAuthTokenString,newCSRFString,err = updateAuthTokenString(oldRefreshTokenString,oldAuthTokenString)

						if err!=nil{
						log.Println("CheckAndRefreshToken->Something went wrong while creating newauth token")
						log.Fatal(err)	
						return;
						}

						newRefreshTokenString,err = updateRefreshTokenExpire(oldRefreshTokenString)

						if err!=nil{
							
						log.Println("CheckAndRefreshToken->Something went wrong while creating new refresh token")
						log.Fatal(err)
							return;
						}
						
						newRefreshTokenString,err = updateRefreshTokenCSRF(newRefreshTokenString,newCSRFString)

						return;
					}else{

						log.Println("CheckAndRefreshToken->Error in auth Token")
						err = errors.New("error in auth token")
						return;
					}
				}else{
					log.Println("error in auth token")
					log.Println("CheckAndRefreshToken->Error in auth Token")
					err = errors.New("error in auth token")
					return;
				}
}