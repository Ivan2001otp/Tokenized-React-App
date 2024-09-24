package middleware

import (
	"Backend/helper"
	"Backend/shared"
	"encoding/json"
	"log"
	"net/http"
)

type status map[string]interface{}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the appropriate CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Must be the frontend origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control,")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		w.Header().Set("Access-Control-Max-Age","3600");
		w.Header().Set("Timing-Allow-Origin","*");

		w.Header().Set("X-Content-Type-Options","nosniff")
		w.Header().Set("X-Frame-Options","DENY")
		w.Header().Set("Referrer-Policy","strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy","default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; base-uri 'self'; form-action 'self'")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials like cookies or tokens

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK) // Respond to OPTIONS with 200 OK
			return
		}

		// Serve the actual request
		next.ServeHTTP(w, r)
	})
}

/*
func CorsMiddleware(next http.Handler) http.Handler{

	return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
		allowedHeaders := " Content-Type,Authorization,X-CSRF-Token"

		w.Header().Set("Access-Control-Allow-Origin","http://localhost:5173");
		w.Header().Set("Access-Control-Allow-Methods","GET,POST,PUT,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",allowedHeaders)
		// w.Header().Set("Access-Control-Expose-Headers","Authorization,X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials","true");

		if(r.Method==http.MethodOptions){
			log.Println("Preflight request made");
			// w.Write([]byte("Preflight Request Made"));
			w.WriteHeader(http.StatusOK)
			return;
		}

		 next.ServeHTTP(w,r);
	})
}*/

// if any crash happens the applicationd handles it gracefully.
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Panicf("Recovered from panic : %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("Entered in auth restricted section")

		log.Println(r)

		//change it to token
		authCookie, authErr := r.Cookie(shared.AUTH_TOKEN)

		log.Println("Auth Token ----> ", authCookie)

		if authErr == http.ErrNoCookie {
			log.Println("Unauthorized attempt! No auth cookie")
			helper.NullifyTokenCookies(&w, r)
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(authErr.Error())
			json.NewEncoder(w).Encode(status{"error": authErr.Error()})
			return
		} else if authErr != nil {
			log.Panicf("panic : %+v", authErr)
			helper.NullifyTokenCookies(&w, r)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": authErr.Error()})
			return
		}

		refreshCookie, refreshErr := r.Cookie(shared.REFRESH_TOKEN)
		log.Println("refres Token ----> ", refreshCookie)

		if refreshErr == http.ErrNoCookie {
			log.Println("Unauthorized attempt! no refresh cookie")
			helper.NullifyTokenCookies(&w, r)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(status{"error": refreshErr.Error()})
			return
		} else if refreshErr != nil {
			log.Panicf("panic:%+v", refreshErr)
			helper.NullifyTokenCookies(&w, r)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"status": refreshErr.Error()})
			return
		}

		//get the csrftoken
		log.Println("here is the csrf token")
		requestCSRFString := helper.GrabCSRFfromRequest(r)
		log.Println("Csrf - ", requestCSRFString)

		if requestCSRFString == "" {
			log.Println("No csrf token passed to request header")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(status{"error": "csrf token is absent in header"})
			return
		}

		authTokenString, refreshTokenString, csrfString, err := helper.CheckAndRefreshToken(authCookie.Value, refreshCookie.Value, requestCSRFString)

		if err != nil {
			if err.Error() == "Unauthorized" {
				log.Println("Unauthorized attempt! JWT not valid!")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(status{"error": "JWT is not valid:"})
				return
			} else {
				log.Println("Error not nil")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(status{"error": "JWT is not valid"})
				return
			}
		}

		log.Println("Successfully recreated jwt")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		helper.SetAuthAndRefreshCookies(&w, authTokenString, refreshTokenString)

		w.Header().Set(shared.X_CSRF_Token, csrfString)

		next.ServeHTTP(w, r)
	})
}

/*
func Authenticator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		log.Println("Trigerred endpoint : ", r.URL.Path)

		switch r.URL.Path {
		case "/dashboard", "/logout", "/deleteUser":
			log.Println("Entered in auth restricted section")

			log.Println(r)

			authCookie, authErr := r.Cookie(shared.AUTH_TOKEN)
			log.Println("Auth Token ----> ", authCookie)

			if authErr == http.ErrNoCookie {
				log.Println("Unauthorized attempt! No auth cookie")
				helper.NullifyTokenCookies(&w, r)
				w.WriteHeader(http.StatusUnauthorized)
				log.Println(authErr.Error())
				json.NewEncoder(w).Encode(status{"error": authErr.Error()})
				return
			} else if authErr != nil {
				log.Panicf("panic : %+v", authErr)
				helper.NullifyTokenCookies(&w, r)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(status{"error": authErr.Error()})
				return
			}

			refreshCookie, refreshErr := r.Cookie(shared.REFRESH_TOKEN)
			log.Println("refres Token ----> ", refreshCookie)

			if refreshErr == http.ErrNoCookie {
				log.Println("Unauthorized attempt! no refresh cookie")
				helper.NullifyTokenCookies(&w, r)
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(status{"error": refreshErr.Error()})
				return
			} else if refreshErr != nil {
				log.Panicf("panic:%+v", refreshErr)
				helper.NullifyTokenCookies(&w, r)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(status{"status": refreshErr.Error()})
				return
			}

			//get the csrftoken
			log.Println("here is the csrf token")
			requestCSRFString := helper.GrabCSRFfromRequest(r)
			log.Println("Csrf - ", requestCSRFString)

			authTokenString, refreshTokenString, csrfString, err := helper.CheckAndRefreshToken(authCookie.Value, refreshCookie.Value, requestCSRFString)

			if err != nil {
				if err.Error() == "Unauthorized" {
					log.Println("Unauthorized attempt! JWT not valid!")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(status{"error": "JWT is not valid:"})
					return
				} else {
					log.Println("Error not nil")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(status{"error": "JWT is not valid"})
					return
				}
			}

			log.Println("Successfully recreated jwt")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			helper.SetAuthAndRefreshCookies(&w, authTokenString, refreshTokenString)

			w.Header().Set(shared.X_CSRF_Token, csrfString)

			break
		default:
*/
