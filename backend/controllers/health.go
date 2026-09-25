package controllers

import (
	"net/http"

	"github.com/IRISBCOSTA/proj-ordert/internal/httpx"
)

func Health(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
