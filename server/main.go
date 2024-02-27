package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const USERNAME = "yasar"
const PASSWORD = "abcde"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(CORS)
	r.Use(middleware.Logger)

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if u.Username == "" || u.Password == "" {
			http.Error(w, "username and password can't be empty", http.StatusUnauthorized)
			return
		}

		if u.Username != USERNAME || u.Password != PASSWORD {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		value := "user-logged-in"

		cookie := &http.Cookie{
			Name:     "session",
			Value:    value,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			MaxAge:   60 * 60 * 24,
		}

		http.SetCookie(w, cookie)

		w.Write([]byte("logged in"))

	})

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("You can view this page because you are logged in"))
		})
	})

	http.ListenAndServe(":8000", r)
}

var CORS = cors.Handler(cors.Options{
	AllowedOrigins: []string{"http://*", "https://*"},
	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
})

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "no cookie", http.StatusUnauthorized)
			return
		}

		if cookie.Value != "user-logged-in" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, nil)
	})
}
