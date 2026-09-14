package routers

import (
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func BuildBudgetRouter() *chi.Mux {
	budgetRouter := chi.NewRouter()

	budgetRouter.Use(middleware.UnifiedAuthMiddleware)
	budgetRouter.Post("/{groupId}", handlers.GetBudgetData)
	budgetRouter.Put("/{groupId}", handlers.UpsertBudget)
	budgetRouter.Delete("/{groupId}/{categoryId}", handlers.DeleteBudget)

	return budgetRouter
}
