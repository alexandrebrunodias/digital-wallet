package web

import (
	"encoding/json"
	"fmt"
	"github.com/alexandrebrunodias/wallet-core/internal/usecase/create_customer"
	"net/http"
)

type CustomerHandler struct {
	CreateCustomerUseCase create_customer.CreateCustomerUseCase
}

func NewCustomerHandler(createCustomerUseCase create_customer.CreateCustomerUseCase) *CustomerHandler {
	if &createCustomerUseCase == nil {
		panic("'CreateCustomerUseCase' must not be nil")
	}
	return &CustomerHandler{CreateCustomerUseCase: createCustomerUseCase}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var command create_customer.CreateCustomerCommand
	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	output, err := h.CreateCustomerUseCase.Execute(command)
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
