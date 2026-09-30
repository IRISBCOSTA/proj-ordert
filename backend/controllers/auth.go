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

// Register cria o usuário E a conta bancária dele (saldo zero) na mesma
// transação
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

	tx, err := database.DB.Begin()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao iniciar transação")
		return
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRow(
		`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		req.Name, req.Email, string(hash),
	).Scan(&userID)
	if err != nil {
		httpx.Error(w, http.StatusConflict, "Email já cadastrado")
		return
	}

	var accountID int
	err = tx.QueryRow(
		`INSERT INTO accounts (user_id, balance) VALUES ($1, 0) RETURNING id`,
		userID,
	).Scan(&accountID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao criar conta bancária")
		return
	}

	if err := tx.Commit(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao confirmar cadastro")
		return
	}

	httpx.JSON(w, http.StatusCreated, map[string]any{
		"id": userID, "name": req.Name, "email": req.Email, "account_id": accountID,
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
