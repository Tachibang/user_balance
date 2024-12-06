package router

import (
	"net/http"
	"user_balance/internal/service"
)

func NewRouter(services *service.Service) *http.ServeMux {
	mux := http.NewServeMux()

	api := "/api"
	newAccountRoutes(mux, api+"/accounts", services.Account)
	newProductRoutes(mux, api+"/products", services.Product)
	newReservationRoutes(mux, api+"/reservations", services.Reservation)

	return mux
}
