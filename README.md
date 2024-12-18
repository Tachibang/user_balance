# User Balance API

## Требования
- Docker
- Docker Compose

## Установка и запуск
1. **.env:**
заполните файл .env

2. **makefile:**
   Выполните следующую команду:

```
make run
```


## Тестирование API

# Запрос на создание аккаунта
curl -X POST http://localhost:8080/api/accounts/create

# Запрос на получение аккаунта
curl -X GET http://localhost:8080/api/accounts/get?id=2

# Запрос на пополнение счета
curl -X POST http://localhost:8080/api/accounts/deposit?id=2&amount=100

# Запрос на снятие средств с аккаунта
curl -X POST http://localhost:8080/api/accounts/withdraw?id=1&amount=100

# Запрос на перевод средств между аккаунтами
curl -X POST http://localhost:8080/api/accounts/transfer?idTo=1&idFrom=2&amount=50

# Запрос на создание продукта
curl -X POST http://localhost:8080/api/products/create?name=cleaning

# Запрос на получение продукта
curl -X GET http://localhost:8080/api/products/get?id=1

# Запрос на создание резервации
curl -X POST http://localhost:8080/api/reservations/create?account_id=1&product_id=1&order_id=1

# Запрос на получение резервации
curl -X GET http://localhost:8080/api/reservations/get?id=2

# Запрос на возврат по резервации
curl -X POST http://localhost:8080/api/reservations/refund?id=2