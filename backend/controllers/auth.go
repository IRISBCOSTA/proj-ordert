package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/IRISBCOSTA/proj-ordert/database"
	"github.com/IRISBCOSTA/proj-ordert/internal/httpx"
	"github.com/IRISBCOSTA/proj-ordert/middleware"
	"github.com/IRISBCOSTA/proj-ordert/models"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if req.Name == "" || req.Email == "" || len(req.Password) < 6 {
		httpx.Error(w, http.StatusBadRequest, "Nome, email e senha (mín. 6 caracteres) são obrigatórios")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao processar senha")
		return
	}

	var id int
	err = database.DB.QueryRow(
		`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		req.Name, req.Email, string(hash),
	).Scan(&id)
	if err != nil {
		httpx.Error(w, http.StatusConflict, "Email já cadastrado")
		return
	}

	httpx.JSON(w, http.StatusCreated, map[string]any{
		"id": id, "name": req.Name, "email": req.Email,
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	var user models.User
	err := database.DB.QueryRow(
		`SELECT id, name, email, password_hash FROM users WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "Email ou senha inválidos")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		httpx.Error(w, http.StatusUnauthorized, "Email ou senha inválidos")
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao gerar token")
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]string{"token": token})
}

// rota protegida: só responde se o
// middleware.RequireAuth já validou o token.
func Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r)

	var user models.User
	err := database.DB.QueryRow(
		`SELECT id, name, email, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Usuário não encontrado")
		return
	}

	httpx.JSON(w, http.StatusOK, user)
}

