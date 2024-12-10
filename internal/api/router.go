package api

import (
	"net/http"
	"user_balance/internal/api/v1/handler"
	"user_balance/internal/service"

	"github.com/sirupsen/logrus"
)

func NewRouter(services *service.Service, logger *logrus.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	apiV1 := "/api/v1"

	mux.HandleFunc(apiV1+"/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.NewAccountRoutes(mux, apiV1+"/accounts", services.Account, logger)
	handler.NewProductRoutes(mux, apiV1+"/products", services.Product, logger)
	handler.NewReservationRoutes(mux, apiV1+"/reservations", services.Reservation, logger)

	return mux
}
