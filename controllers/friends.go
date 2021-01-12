package controllers

import (
	"fmt"
	"net/http"
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
	// user := context.User(r.Context())
	// friends, err := f.fs.ByUserID(user.ID)

	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	// 	return
	// }
	// var vd views.Data
	// vd.Yield = friends
	// g.IndexView.Render(w, r, vd)
}

// GET /galleries/:id
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
	fmt.Println("create: ", r)
}
