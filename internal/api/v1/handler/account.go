package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"user_balance/internal/service"

	"github.com/sirupsen/logrus"
)

func NewAccountRoutes(mux *http.ServeMux, basePath string, accountService service.Account, logger *logrus.Logger) {
	mux.HandleFunc(basePath+"/create", createAccountHandler(accountService, logger))
	mux.HandleFunc(basePath+"/get", getAccountHandler(accountService, logger))
	mux.HandleFunc(basePath+"/deposit", depositAccountHandler(accountService, logger))
	mux.HandleFunc(basePath+"/withdraw", withdrawAccountHandler(accountService, logger))
	mux.HandleFunc(basePath+"/transfer", transferAccountHandler(accountService, logger))
}

func createAccountHandler(accountService service.Account, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.Warnf("Запрос к %s не выполнен: метод не разрешён", r.URL.Path)
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		id, err := accountService.CreateAccount(r.Context())
		if err != nil {
			logger.Errorf("Не удалось создать аккаунт: %v", err)
			http.Error(w, "Не удалось создать аккаунт", http.StatusInternalServerError)
			return
		}
		type response struct {
			Id int `json:"id"`
		}

		logger.Infof("Аккаунт успешно создан с ID %d", id)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response{
			Id: id,
		})
	}
}

func getAccountHandler(accountService service.Account, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			logger.Warnf("Запрос к %s не выполнен: метод не разрешён", r.URL.Path)
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			logger.Warnf("Запрос к %s не выполнен: отсутствует или недопустим ID аккаунта", r.URL.Path)
			http.Error(w, "Отсутствует или недопустим ID аккаунта", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Warnf("Запрос к %s не выполнен: недопустимый формат ID аккаунта", r.URL.Path)
			http.Error(w, "Недопустимый формат ID аккаунта", http.StatusBadRequest)
			return
		}

		account, err := accountService.GetAccount(r.Context(), id)
		if err != nil {
			logger.Errorf("Не удалось получить аккаунт с ID %d: %v", id, err)
			http.Error(w, "Не удалось получить аккаунт", http.StatusInternalServerError)
			return
		}

		type response struct {
			Id      int `json:"id"`
			Balance int `json:"balance"`
		}

		logger.Infof("Аккаунт успешно получен: ID %d, Баланс %d", account.Id, account.Balance)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{
			Id:      account.Id,
			Balance: account.Balance,
		})
	}
}

func depositAccountHandler(accountService service.Account, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.Warnf("Запрос к %s не выполнен: метод не разрешён", r.URL.Path)
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		idParam := r.URL.Query().Get("id")
		amountParam := r.URL.Query().Get("amount")

		id, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Warnf("Запрос к %s не выполнен: недопустимый формат ID аккаунта", r.URL.Path)
			http.Error(w, "Недопустимый формат ID аккаунта", http.StatusBadRequest)
			return
		}

		amount, err := strconv.Atoi(amountParam)
		if err != nil {
			logger.Warnf("Запрос к %s не выполнен: недопустимый формат суммы", r.URL.Path)
			http.Error(w, "Недопустимый формат суммы", http.StatusBadRequest)
			return
		}

		updatedId, updatedBalance, err := accountService.Deposit(r.Context(), id, amount)
		if err != nil {
			logger.Errorf("Не удалось обновить баланс аккаунта с ID %d: %v", id, err)
			http.Error(w, "Не удалось обновить баланс аккаунта", http.StatusInternalServerError)
			return
		}

		logger.Infof("Пополнение успешно: ID %d, Новый баланс %d", updatedId, updatedBalance)
		w.WriteHeader(http.StatusOK)
		type response struct {
			Id      int `json:"id"`
			Balance int `json:"balance"`
		}
		json.NewEncoder(w).Encode(response{
			Id:      updatedId,
			Balance: updatedBalance,
		})
	}
}

func withdrawAccountHandler(accountService service.Account, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.Warnf("Запрос к %s не выполнен: метод не разрешён", r.URL.Path)
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		idParam := r.URL.Query().Get("id")
		amountParam := r.URL.Query().Get("amount")

		id, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Warnf("Запрос к %s не выполнен: недопустимый формат ID аккаунта", r.URL.Path)
			http.Error(w, "Недопустимый формат ID аккаунта", http.StatusBadRequest)
			return
		}

		amount, err := strconv.Atoi(amountParam)
		if err != nil {
			logger.Warnf("Запрос к %s не выполнен: недопустимый формат суммы", r.URL.Path)
			http.Error(w, "Недопустимый формат суммы", http.StatusBadRequest)
			return
		}

		updatedId, updatedBalance, err := accountService.Withdraw(r.Context(), id, amount)
		if err != nil {
			logger.Errorf("Не удалось обновить баланс аккаунта с ID %d: %v", id, err)
			http.Error(w, "Не удалось обновить баланс аккаунта", http.StatusInternalServerError)
			return
		}

		logger.Infof("Снятие успешно: ID %d, Новый баланс %d", updatedId, updatedBalance)
		w.WriteHeader(http.StatusOK)
		type response struct {
			Id      int `json:"id"`
			Balance int `json:"balance"`
		}
		json.NewEncoder(w).Encode(response{
			Id:      updatedId,
			Balance: updatedBalance,
		})
	}
}

func transferAccountHandler(accountService service.Account, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.Warnf("Запрос к %s не выполнен: метод не разрешён", r.URL.Path)
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		idToParam := r.URL.Query().Get("idTo")
		idFromParam := r.URL.Query().Get("idFrom")
		amountParam := r.URL.Query().Get("amount")

		idTo, err := strconv.Atoi(idToParam)
		if err != nil {
			logger.Warnf("Запрос к %s не выполнен: недопустимый формат ID аккаунта (idTo)", r.URL.Path)
			http.Error(w, "Недопустимый формат ID аккаунта", http.StatusBadRequest)
			return
		}

		idFrom, err := strconv.Atoi(idFromParam)
		if err != nil {
			logger.Warnf("Запрос к %s не выполнен: недопустимый формат ID аккаунта (idFrom)", r.URL.Path)
			http.Error(w, "Недопустимый формат ID аккаунта", http.StatusBadRequest)
			return
		}

		amount, err := strconv.Atoi(amountParam)
		if err != nil {
			logger.Warnf("Запрос к %s не выполнен: недопустимый формат суммы", r.URL.Path)
			http.Error(w, "Недопустимый формат суммы", http.StatusBadRequest)
			return
		}

		updatedBalanceTo, updatedBalanceFrom, err := accountService.Transfer(r.Context(), idTo, idFrom, amount)
		if err != nil {
			logger.Errorf("Не удалось обновить баланс при переводе с ID %d на ID %d: %v", idFrom, idTo, err)
			http.Error(w, "Не удалось обновить баланс при переводе", http.StatusInternalServerError)
			return
		}

		logger.Infof("Перевод успешно выполнен: С ID %d, на ID %d, Сумма %d", idFrom, idTo, amount)
		type response struct {
			BalanceTo   int `json:"to"`
			BalanceFrom int `json:"from"`
			Amount      int `json:"amount"`
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{
			BalanceTo:   updatedBalanceTo,
			BalanceFrom: updatedBalanceFrom,
			Amount:      amount,
		})
	}
}
