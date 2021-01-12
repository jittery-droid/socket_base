package middleware

import (
	"fmt"
	"net/http"
	"sockets/context"
	"sockets/models"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	models.UserService
}

func (mw *User) extractToken(r *http.Request) string {
	keys := r.URL.Query()
	fmt.Println("keys:", keys)
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	fmt.Println(bearerToken)
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func (mw *User) extractUser(tokenString string) uint {
	fmt.Println("tokenString", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing algo")
		}
		return []byte(mw.JwtSecret()), nil
	})
	if err != nil {
		panic(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			panic(err)
		}
		return uint(uid)
	}
	return 0
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			next(w, r)
			return
		}

		path := r.URL.Path

		if strings.HasPrefix(path, "/api/login") ||
			strings.HasPrefix(path, "/api/register") ||
			strings.HasPrefix(path, "/api/auth") ||
			strings.HasPrefix(path, "/") {
			next(w, r)
			return
		}

		tokenString := mw.extractToken(r)
		userID := mw.extractUser(tokenString)
		user, err := mw.UserService.ByID(userID)
		if err != nil {
			next(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

// RequireUser assumes that User middleware has already been run
// otherwise it will no work correctly.
type RequireUser struct {
	User
}

// Apply assumes that User middleware has already been run
// otherwise it will no work correctly.
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn assumes that User middleware has already been run
// otherwise it will no work correctly.
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
