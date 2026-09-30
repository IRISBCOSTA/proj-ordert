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

	mux.HandleFunc("GET /account", middleware.RequireAuth(controllers.GetAccount))
	mux.HandleFunc("POST /deposit", middleware.RequireAuth(controllers.Deposit))
	mux.HandleFunc("POST /withdraw", middleware.RequireAuth(controllers.Withdraw))
	mux.HandleFunc("POST /transfer", middleware.RequireAuth(controllers.Transfer))
	mux.HandleFunc("GET /transactions", middleware.RequireAuth(controllers.GetTransactions))

	return mux
}
