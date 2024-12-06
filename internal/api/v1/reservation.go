package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"user_balance/internal/entity"
	"user_balance/internal/service"
)

func newReservationRoutes(mux *http.ServeMux, basePath string, reservationService service.Reservation) {
	mux.HandleFunc(basePath+"/create", createReservationHandler(reservationService))
	mux.HandleFunc(basePath+"/get", getReservationHandler(reservationService))
	mux.HandleFunc(basePath+"/refund", refundReservationHandler(reservationService))
}

func createReservationHandler(reservationService service.Reservation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			log.Println("Ошибка: попытка использования недопустимого метода для создания резервации")
			return
		}

		accountIDParam := r.URL.Query().Get("account_id")
		productIDParam := r.URL.Query().Get("product_id")
		amountParam := r.URL.Query().Get("amount")

		if accountIDParam == "" || productIDParam == "" || amountParam == "" {
			http.Error(w, "Отсутствуют или некорректные параметры", http.StatusBadRequest)
			log.Println("Ошибка: отсутствуют или некорректные параметры запроса")
			return
		}

		accountID, err := strconv.Atoi(accountIDParam)
		if err != nil {
			http.Error(w, "Неверный формат ID аккаунта", http.StatusBadRequest)
			log.Printf("Ошибка: неверный формат ID аккаунта: %v\n", err)
			return
		}

		productID, err := strconv.Atoi(productIDParam)
		if err != nil {
			http.Error(w, "Неверный формат ID продукта", http.StatusBadRequest)
			log.Printf("Ошибка: неверный формат ID продукта: %v\n", err)
			return
		}

		amount, err := strconv.Atoi(amountParam)
		if err != nil {
			http.Error(w, "Неверный формат суммы", http.StatusBadRequest)
			log.Printf("Ошибка: неверный формат суммы: %v\n", err)
			return
		}

		reservation := entity.Reservation{
			AccountId: accountID,
			ProductId: productID,
			Amount:    amount,
		}

		reservationID, err := reservationService.CreateReservation(r.Context(), reservation)
		if err != nil {
			http.Error(w, "Не удалось создать резервацию", http.StatusInternalServerError)
			log.Printf("Ошибка при создании резервации: %v\n", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(reservationID)
		log.Printf("Резервация с ID %d успешно создана\n", reservationID)
	}
}

func getReservationHandler(reservationService service.Reservation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			log.Println("Ошибка: попытка использования недопустимого метода для получения резервации")
			return
		}

		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Отсутствует или некорректный ID резервации", http.StatusBadRequest)
			log.Println("Ошибка: отсутствует или некорректный ID резервации")
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Неверный формат ID резервации", http.StatusBadRequest)
			log.Printf("Ошибка: неверный формат ID резервации: %v\n", err)
			return
		}

		reservation, err := reservationService.GetReservation(r.Context(), id)
		if err != nil {
			http.Error(w, "Не удалось получить резервацию", http.StatusInternalServerError)
			log.Printf("Ошибка при получении резервации с ID %d: %v\n", id, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(reservation)
		log.Printf("Резервация с ID %d успешно получена\n", id)
	}
}

func refundReservationHandler(reservationService service.Reservation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			log.Println("Ошибка: попытка использования недопустимого метода для возврата резервации")
			return
		}

		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Отсутствует или некорректный ID резервации", http.StatusBadRequest)
			log.Println("Ошибка: отсутствует или некорректный ID резервации")
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Неверный формат ID резервации", http.StatusBadRequest)
			log.Printf("Ошибка: неверный формат ID резервации: %v\n", err)
			return
		}

		err = reservationService.RefundReservation(r.Context(), id)
		if err != nil {
			http.Error(w, "Не удалось вернуть резервацию", http.StatusInternalServerError)
			log.Printf("Ошибка при возврате резервации с ID %d: %v\n", id, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Резервация успешно возвращена")
		log.Printf("Резервация с ID %d успешно возвращена\n", id)
	}
}
