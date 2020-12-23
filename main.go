package main

import (
	"flag"
	"fmt"
	"net/http"
	"sockets/controllers"
	"sockets/middleware"
	"sockets/models"
	"sockets/rand"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the application starts.")
	flag.Parse()

	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
	)
	must(err)
	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()
	usersC := controllers.NewUsers(services.User)

	b, err := rand.Bytes(32)
	must(err)
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")

	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
