package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IRISBCOSTA/proj-ordert/database"
	"github.com/IRISBCOSTA/proj-ordert/internal/httpx"
	"github.com/IRISBCOSTA/proj-ordert/middleware"
	"github.com/IRISBCOSTA/proj-ordert/models"
)

//busca conta
func getAccountByUserID(userID int) (models.Account, error) {
	var account models.Account
	err := database.DB.QueryRow(
		`SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = $1`,
		userID,
	).Scan(&account.ID, &account.UserID, &account.Balance, &account.CreatedAt)
	return account, err
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	account, err := getAccountByUserID(middleware.UserID(r))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Conta não encontrada")
		return
	}

	httpx.JSON(w, http.StatusOK, account)
}

type amountRequest struct {
	Amount float64 `json:"amount"`
}

func Deposit(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r)

	var req amountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		httpx.Error(w, http.StatusBadRequest, "Valor de depósito inválido")
		return
	}

	account, err := getAccountByUserID(userID)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Conta não encontrada")
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao iniciar transação")
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, req.Amount, account.ID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao atualizar saldo")
		return
	}

	if _, err := tx.Exec(
		`INSERT INTO transactions (account_id, type, amount) VALUES ($1, 'deposito', $2)`,
		account.ID, req.Amount,
	); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao registrar depósito")
		return
	}

	if err := tx.Commit(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao confirmar depósito")
		return
	}

	updated, _ := getAccountByUserID(userID)
	httpx.JSON(w, http.StatusOK, updated)
}

func Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r)

	var req amountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		httpx.Error(w, http.StatusBadRequest, "Valor de saque inválido")
		return
	}

	account, err := getAccountByUserID(userID)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Conta não encontrada")
		return
	}

	if account.Balance < req.Amount {
		httpx.Error(w, http.StatusBadRequest, "Saldo insuficiente")
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao iniciar transação")
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, req.Amount, account.ID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao atualizar saldo")
		return
	}

	if _, err := tx.Exec(
		`INSERT INTO transactions (account_id, type, amount) VALUES ($1, 'saque', $2)`,
		account.ID, req.Amount,
	); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao registrar saque")
		return
	}

	if err := tx.Commit(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao confirmar saque")
		return
	}

	updated, _ := getAccountByUserID(userID)
	httpx.JSON(w, http.StatusOK, updated)
}

type statementEntry struct {
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	account, err := getAccountByUserID(middleware.UserID(r))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Conta não encontrada")
		return
	}

	rows, err := database.DB.Query(`
		SELECT type, amount, created_at FROM transactions WHERE account_id = $1
		UNION ALL
		SELECT 'transferencia_enviada', amount, created_at FROM transfers WHERE from_account = $1
		UNION ALL
		SELECT 'transferencia_recebida', amount, created_at FROM transfers WHERE to_account = $1
		ORDER BY created_at DESC
	`, account.ID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao buscar extrato")
		return
	}
	defer rows.Close()

	entries := []statementEntry{}
	for rows.Next() {
		var e statementEntry
		if err := rows.Scan(&e.Type, &e.Amount, &e.CreatedAt); err != nil {
			httpx.Error(w, http.StatusInternalServerError, "Erro ao ler extrato")
			return
		}
		entries = append(entries, e)
	}

	httpx.JSON(w, http.StatusOK, entries)
}
