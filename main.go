package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sockets/controllers"
	"sockets/middleware"
	"sockets/models"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP serves static js assets
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	fmt.Println("In spa handler")
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(h.staticPath, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	csrf.TemplateField(r)
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the application starts.")
	flag.Parse()

	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey, cfg.JWTSecret),
	)
	must(err)
	defer services.Close()
	services.AutoMigrate()

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r := mux.NewRouter()
	usersC := controllers.NewUsers(services.User)

	r.HandleFunc("/api/auth", usersC.Load).Methods("GET")
	r.HandleFunc("/api/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/api/login", usersC.Login).Methods("POST")
	r.HandleFunc("/api/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")

	spa := spaHandler{staticPath: "client/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	srv := &http.Server{
		Handler:      userMw.Apply(r),
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
