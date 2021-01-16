package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sockets/context"
	"sockets/models"

	"github.com/gorilla/mux"
)

type Friends struct {
	fs models.FriendService
	r  *mux.Router
}

func NewFriends(fs models.FriendService, r *mux.Router) *Friends {
	return &Friends{
		fs: fs,
		r:  r,
	}
}

// GET /friends
func (f *Friends) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	friends, err := f.fs.ByUserID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(friends)
	json.NewEncoder(w).Encode(friends)
}

func (f *Friends) Show(w http.ResponseWriter, r *http.Request) {
	// gallery, err := g.galleryByID(w, r)
	// if err != nil {
	// 	return
	// }
	// var vd views.Data
	// vd.Yield = gallery
	// g.ShowView.Render(w, r, vd)
}

func (f *Friends) Create(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	var friend models.Friend
	json.NewDecoder(r.Body).Decode(&friend)
	friend.Status = "pending"
	friend.UserID = user.ID
	err := f.fs.Create(&friend)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}
