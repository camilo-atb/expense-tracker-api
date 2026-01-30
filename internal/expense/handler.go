package expense

import (
	"encoding/json"
	"errors"
	"expense-tracker/internal/shared/httpx"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type Handler struct {
	service Service
}

func NewHandler(serviceIny Service) *Handler {
	return &Handler{service: serviceIny}
}

func (h *Handler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var create CreateDTO
	err := json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.service.AddTransaction(&create)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.ResponseHandler(w, http.StatusCreated, transaction)
}

func (h *Handler) ModifyTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var update UpdateDTO

	err = json.NewDecoder(r.Body).Decode(&update)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.service.ModifyTransaction(id, &update)

	if err != nil {
		writeError(w, err)
		return
	}

	httpx.ResponseHandler(w, http.StatusOK, transaction)
}

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	mode := r.URL.Query().Get("mode")

	switch mode {
	case "soft":
		err = h.service.DeleteLogicalTransaction(id)
	case "hard":
		err = h.service.DeleteTransaction(id)
	default:
		http.Error(w, "invalid delete mode", http.StatusBadRequest)
		return
	}
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.service.GetTransactionById(id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.ResponseHandler(w, http.StatusOK, transaction)
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	txType := strings.TrimSpace(q.Get("type"))
	from := strings.TrimSpace(q.Get("from"))
	to := strings.TrimSpace(q.Get("to"))

	var (
		transactions any
		err          error
	)

	switch {
	case txType != "" && from != "" && to != "":
		transactions, err = h.service.GetTransactionsByTypeAndDate(TransactionType(txType), from, to)

	case from != "" && to != "":
		transactions, err = h.service.GetTransactionsBetweenDates(from, to)

	case txType != "":
		transactions, err = h.service.GetTransactionsByType(TransactionType(txType))

	default:
		transactions, err = h.service.GetAllTransactions()
	}

	if err != nil {
		writeError(w, err)
		return
	}

	httpx.ResponseHandler(w, http.StatusOK, transactions)
}

func (h *Handler) GetTotalsByType(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	txType := strings.TrimSpace(q.Get("type"))
	from := strings.TrimSpace(q.Get("from"))
	to := strings.TrimSpace(q.Get("to"))

	switch TransactionType(txType) {
	case Income, Expense:
		total, err := h.service.ObtainTotalByTypeAndDate(TransactionType(txType), from, to)
		if err != nil {
			writeError(w, err)
			return
		}
		httpx.ResponseHandler(w, http.StatusOK, total)
	default:
		http.Error(w, "invalid type", http.StatusBadRequest)
	}

}

func (h *Handler) GetNetIncome(w http.ResponseWriter, r *http.Request) {
	total, err := h.service.GetNetIncome()
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.ResponseHandler(w, http.StatusOK, total)
}

// --------

func getId(request *http.Request) (int, error) {
	idStr := chi.URLParam(request, "id")
	if idStr == "" {
		return 0, errors.New("missing id")
	}
	return strconv.Atoi(idStr)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, errInvalidId):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrTransactionNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
