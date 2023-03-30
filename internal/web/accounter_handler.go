package web

import (
	"encoding/json"
	"fmt"
	"github.com/alexandrebrunodias/wallet-core/internal/usecase/create_account"
	"net/http"
)

type AccountHandler struct {
	CreateAccountUseCase create_account.CreateAccountUseCase
}

func NewAccountHandler(createAccountUseCase create_account.CreateAccountUseCase) *AccountHandler {
	if &createAccountUseCase == nil {
		panic("'CreateAccountUseCase' must not be nil")
	}
	return &AccountHandler{CreateAccountUseCase: createAccountUseCase}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var command create_account.CreateAccountCommand
	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	output, err := h.CreateAccountUseCase.Execute(command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}
