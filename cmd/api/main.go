package main

import (
	"errors"
	"expense-tracker/internal/expense"
	"expense-tracker/internal/shared/database"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env:", err)
	}
	db, err := database.DatabaseConnection()
	if err != nil {
		log.Fatal("Error en la base de datos:", err)
	}
	repo := expense.NewRepository(db)
	service := expense.NewService(repo)
	handler := expense.NewHandler(service)

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {

		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", handler.GetTransactions)
			r.Get("/{id}", handler.GetTransactionById)
			r.Post("/", handler.AddTransaction)
			r.Patch("/{id}", handler.ModifyTransaction)
			r.Delete("/{id}", handler.DeleteTransaction)
		})

		r.Route("/summary", func(r chi.Router) {
			r.Get("/type", handler.GetTotalsByType)
			r.Get("/net", handler.GetNetIncome)
		})
	})

	startServer(r)

}

func startServer(handler http.Handler) {
	addr := ":8080"

	log.Printf("API listening on %s\n", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
