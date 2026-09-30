package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/IRISBCOSTA/proj-ordert/database"
	"github.com/IRISBCOSTA/proj-ordert/internal/httpx"
	"github.com/IRISBCOSTA/proj-ordert/middleware"
)

type transferRequest struct {
	ToAccountID int     `json:"to_account_id"`
	Amount      float64 `json:"amount"`
}

// Transfer move dinheiro entre duas contas na mesma transação: se qualquer
// passo falhar, o defer tx.Rollback() desfaz tudo
func Transfer(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r)

	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		httpx.Error(w, http.StatusBadRequest, "Dados de transferência inválidos")
		return
	}

	fromAccount, err := getAccountByUserID(userID)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Conta de origem não encontrada")
		return
	}

	if req.ToAccountID == fromAccount.ID {
		httpx.Error(w, http.StatusBadRequest, "Não é possível transferir para a própria conta")
		return
	}

	if fromAccount.Balance < req.Amount {
		httpx.Error(w, http.StatusBadRequest, "Saldo insuficiente")
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao iniciar transação")
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, req.Amount, fromAccount.ID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao debitar conta de origem")
		return
	}

	result, err := tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, req.Amount, req.ToAccountID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao creditar conta de destino")
		return
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		httpx.Error(w, http.StatusBadRequest, "Conta de destino não encontrada")
		return
	}

	if _, err := tx.Exec(
		`INSERT INTO transfers (from_account, to_account, amount) VALUES ($1, $2, $3)`,
		fromAccount.ID, req.ToAccountID, req.Amount,
	); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao registrar transferência")
		return
	}

	if err := tx.Commit(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Erro ao confirmar transferência")
		return
	}

	updated, _ := getAccountByUserID(userID)
	httpx.JSON(w, http.StatusOK, updated)
}
