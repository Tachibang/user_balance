package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"user_balance/internal/service"

	"github.com/sirupsen/logrus"
)

func NewProductRoutes(mux *http.ServeMux, basePath string, productService service.Product, logger *logrus.Logger) {
	mux.HandleFunc(basePath+"/create", createProductHandler(productService))
	mux.HandleFunc(basePath+"/get", getProductHandler(productService))
}

func createProductHandler(productService service.Product) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			log.Println("Ошибка: попытка использования недопустимого метода для создания продукта")
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Отсутствует или некорректное имя продукта", http.StatusBadRequest)
			log.Println("Ошибка: отсутствие или некорректное имя продукта")
			return
		}

		id, err := productService.CreateProduct(r.Context(), name)
		if err != nil {
			http.Error(w, "Не удалось создать продукт", http.StatusInternalServerError)
			log.Printf("Ошибка при создании продукта: %v\n", err)
			return
		}

		type response struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response{
			Id:   id,
			Name: name,
		})
		log.Printf("Продукт с ID %d успешно создан\n", id)
	}
}

func getProductHandler(productService service.Product) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			log.Println("Ошибка: попытка использования недопустимого метода для получения продукта")
			return
		}

		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Отсутствует или некорректный ID продукта", http.StatusBadRequest)
			log.Println("Ошибка: отсутствие или некорректный ID продукта")
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Неверный формат ID продукта", http.StatusBadRequest)
			log.Printf("Ошибка: неверный формат ID продукта: %v\n", err)
			return
		}

		product, err := productService.GetProduct(r.Context(), id)
		if err != nil {
			http.Error(w, "Не удалось получить продукт", http.StatusInternalServerError)
			log.Printf("Ошибка при получении продукта с ID %d: %v\n", id, err)
			return
		}

		type response struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{
			Id:   product.Id,
			Name: product.Name,
		})
		log.Printf("Продукт с ID %d успешно получен\n", id)
	}
}
