package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/karthikbhandary2/inventory/internal/db"
	"github.com/karthikbhandary2/inventory/internal/handlers"
	"github.com/karthikbhandary2/inventory/internal/middleware"
	"github.com/karthikbhandary2/inventory/internal/service"
)

func main() {
	ctx := context.Background()

	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	// Wire up dependencies
	productSvc := &service.ProductService{}
	productHandler := handlers.NewProductHandler(productSvc)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(os.Getenv("JWT_SECRET"))) // require valid JWT
		r.Use(middleware.Transaction(pool))             // wrap in DB tx

		// Products
		r.Get("/products", productHandler.List)
		r.Post("/products", productHandler.Create)
		r.Get("/products/{id}", productHandler.Get)
		r.Put("/products/{id}", productHandler.Update)
		r.Delete("/products/{id}", productHandler.Delete)

		// Stock operations (sub-resource)
		r.Post("/products/{id}/stock", productHandler.StockOp)

		// Reports
		r.Get("/reports/inventory", productHandler.Report)

		// Audit log
		r.Get("/audit", productHandler.AuditLog)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
