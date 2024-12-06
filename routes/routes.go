package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	controllers "github.com/jofan-cah/login-api/controller"
	"github.com/jofan-cah/login-api/middlewares"
)

func RegisterRoutes(r *mux.Router) {
	// Register rute untuk registrasi dan login tanpa token
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	// Rute /users akan dilindungi oleh middleware VerifyToken
	r.Handle("/users", middlewares.VerifyToken(http.HandlerFunc(controllers.GetAllUsers))).Methods("GET")

	// Rute delete user
	r.Handle("/user/{id}", middlewares.VerifyToken(http.HandlerFunc(controllers.DeleteUser))).Methods("DELETE")

	// Rute root untuk testing
	// Rute root untuk testing
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
}
