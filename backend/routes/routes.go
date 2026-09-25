package routes

import (
	"net/http"

	"github.com/IRISBCOSTA/proj-ordert/controllers"
	"github.com/IRISBCOSTA/proj-ordert/middleware"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", controllers.Health)
	mux.HandleFunc("POST /register", controllers.Register)
	mux.HandleFunc("POST /login", controllers.Login)
	mux.HandleFunc("GET /me", middleware.RequireAuth(controllers.Me))

	return mux
}
