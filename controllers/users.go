package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sockets/context"
	"sockets/models"
	"sockets/rand"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Users struct {
	us models.UserService
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(us models.UserService) *Users {
	return &Users{
		us: us,
	}
}

// Load decodes a JWT token and returns the user
func (u *Users) Load(w http.ResponseWriter, r *http.Request) {
	// ADD THIS TO MIDDLEWARE
	tokenString := u.extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing algo")
		}
		return []byte(u.us.JwtSecret()), nil
	})
	if err != nil {
		panic(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		b, err := json.MarshalIndent(claims, "", " ")
		if err != nil {
			panic(err)
		}
		fmt.Println(b)
	}
}

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	if err := u.us.Create(&user); err != nil {
		panic(err)
	}
	jwtToken, err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	err = json.NewEncoder(w).Encode(jwtToken)
	if err != nil {
		panic(err)
	}
}

// Login is used to verify the provided email address and
// password and then log the user in if they are correct.
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:

		default:

		}
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	jwtToken, err := u.signIn(w, user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	err = json.NewEncoder(w).Encode(jwtToken)
	if err != nil {
		panic(err)
	}
}

// Logout is used to delete a users session cookie (remember_token)
// and then will update the user resource with a new remmeber
// token.
//
// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)
}

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) (string, error) {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return "", err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return "", err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	token, err := u.createToken(user)
	if err != nil {
		panic(err)
	}

	return token, nil
}

func (u *Users) createToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.us.JwtSecret()))
}

func (u *Users) extractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
