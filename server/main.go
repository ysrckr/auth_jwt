package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/chmike/securecookie"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(CORS)
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := CreateCookie(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cookie.SetValue(w, []byte("Hello World!")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		val, err := cookie.GetValue(nil, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(html.EscapeString(string(val)))
		w.Write([]byte("Hiiii"))
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

var key []byte = securecookie.MustGenerateRandomKey()

func CreateCookie(key []byte) (*securecookie.Obj, error) {
	cookie, err := securecookie.New("session", key, securecookie.Params{
		Path:     "/",
		MaxAge:   3600,
		HTTPOnly: true,
		Secure:   false,
		SameSite: securecookie.None,
	})
	if err != nil {
		return nil, err
	}
	return cookie, nil
}
