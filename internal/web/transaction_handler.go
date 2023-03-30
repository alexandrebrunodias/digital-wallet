package web

import (
	"encoding/json"
	"fmt"
	"github.com/alexandrebrunodias/wallet-core/internal/usecase/create_transaction"
	"net/http"
)

type TransactionHandler struct {
	CreateTransactionUseCase create_transaction.CreateTransactionUseCase
}

func NewTransactionHandler(createTransactionUseCase create_transaction.CreateTransactionUseCase) *TransactionHandler {
	if &createTransactionUseCase == nil {
		panic("'CreateTransactionUseCase' must not be nil")
	}
	return &TransactionHandler{CreateTransactionUseCase: createTransactionUseCase}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var command create_transaction.CreateTransactionCommand
	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// TODO IMPROVE ERROR HANDLING
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		fmt.Println(err)
		return
	}

	requestContext := r.Context()
	output, err := h.CreateTransactionUseCase.Execute(requestContext, command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}
