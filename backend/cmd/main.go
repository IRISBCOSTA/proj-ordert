package main

import (
	"log"
	"net/http"
	"os"

	"github.com/IRISBCOSTA/proj-ordert/database"
	"github.com/IRISBCOSTA/proj-ordert/middleware"
	"github.com/IRISBCOSTA/proj-ordert/routes"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer database.DB.Close()

	mux := routes.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Ordet Bank API rodando em http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, middleware.CORS(mux)); err != nil {
		log.Fatal(err)
	}
}

