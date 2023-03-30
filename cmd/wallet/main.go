package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexandrebrunodias/wallet-core/internal/database/postgres"
	"github.com/alexandrebrunodias/wallet-core/internal/usecase/create_account"
	"github.com/alexandrebrunodias/wallet-core/internal/usecase/create_customer"
	"github.com/alexandrebrunodias/wallet-core/internal/usecase/create_transaction"
	"github.com/alexandrebrunodias/wallet-core/internal/web"
	"github.com/alexandrebrunodias/wallet-core/pkg/events"
	"github.com/alexandrebrunodias/wallet-core/pkg/events/kafka"
	"github.com/alexandrebrunodias/wallet-core/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

const TransactionsTopic = "wallet.transactions"

func main() {
	// TODO IMPROVE CONFIG
	pgHost := os.Getenv("POSTGRES_HOST")

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			"postgres",
			"senha",
			pgHost,
			"5432",
			"wallet",
		),
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	customerGateway := postgres.NewCustomerPgGateway(db)
	accountGateway := postgres.NewAccountPgGateway(db)

	ctx := context.Background()
	unitOfWork := uow.NewUnitOfWork(ctx, db)
	unitOfWork.Add("AccountGateway", func(tx *sql.Tx) interface{} {
		return postgres.NewAccountPgGateway(db)
	})
	unitOfWork.Add("TransactionGateway", func(tx *sql.Tx) interface{} {
		return postgres.NewTransactionPgGateway(db)
	})

	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", kafkaHost, kafkaPort),
	}

	kafkaProducer := kafka.NewKafkaProducer(configMap, TransactionsTopic, nil)
	transactionEventPublisher := events.NewEventPublisher(kafkaProducer)

	createCustomerUseCase := create_customer.NewCreateCustomerUseCase(customerGateway)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountGateway, customerGateway)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(unitOfWork, transactionEventPublisher)

	customerHandler := web.NewCustomerHandler(*createCustomerUseCase)
	accountHandler := web.NewAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewTransactionHandler(*createTransactionUseCase)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/customers", customerHandler.CreateCustomer)
	router.Post("/accounts", accountHandler.CreateAccount)
	router.Post("/transactions", transactionHandler.CreateTransaction)

	webServerPort := ":8000"
	fmt.Println("Server is running on port", webServerPort)
	panic(http.ListenAndServe(webServerPort, router))
}
