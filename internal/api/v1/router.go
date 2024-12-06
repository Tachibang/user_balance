package router

import (
	"net/http"
	"user_balance/internal/service"
)

func NewV1Router(services *service.Service) *http.ServeMux {
	mux := http.NewServeMux()

	apiV1 := "/api/v1"
	newAccountRoutes(mux, apiV1+"/accounts", services.Account)
	newProductRoutes(mux, apiV1+"/products", services.Product)
	newReservationRoutes(mux, apiV1+"/reservations", services.Reservation)

	// health check endpoint
	mux.HandleFunc(apiV1+"/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}
