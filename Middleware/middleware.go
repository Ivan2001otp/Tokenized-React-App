package middleware

import (
	"log"
	"net/http"
	"time"
)

func SetAuthandRefreshCookies(w *http.ResponseWriter,authTokenString string,refreshTokenString string){
	//auth cookie
	authCookie := http.Cookie{
		Name:"AuthToken",
		Value: authTokenString,
		Secure: true,
		HttpOnly: true,
	}

	authCookie.SameSite = http.SameSiteLaxMode

	http.SetCookie(*w,&authCookie);

	//refresh cookie
	refreshCookie := http.Cookie{
		Name:"RefreshToken",
		Value: refreshTokenString,
		Secure:true,
		HttpOnly: true,
	}

	refreshCookie.SameSite = http.SameSiteLaxMode;

	http.SetCookie(*w,&refreshCookie)
}

func NullifyTokenCookies(w *http.ResponseWriter,r *http.Request){
	authCookie := http.Cookie{
		Name: "AuthToken",
		Value: "",
		Expires: time.Now().Add(-1000*time.Hour),
		HttpOnly:true,

	}

	http.SetCookie(*w,&authCookie)

	refreshCookie := http.Cookie{
		Name: "RefreshToken",
		Value: "",
		Expires: time.Now().Add(-1000*time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(*w,&refreshCookie);

	RefreshCookie,refreshErr := r.Cookie("RefreshToken")

	if refreshErr==http.ErrNoCookie{
		log.Println("No refresh token!")
		return;
	}else if refreshErr!=nil{
		log.Panic("Panic while nullifying tokens : %+v",refreshErr)
		http.Error(*w,http.StatusText(500),500);
	}

	log.Println(RefreshCookie.Value)
	//revoke 
}

func FetchCsrfFromRequest(r *http.Request) string{
	csrfFromRequest := r.FormValue("X-CSRF-Token")

	if csrfFromRequest!=""{
		return csrfFromRequest
	}else{
		return r.Header.Get("X-CSRF-Token")
	}
}