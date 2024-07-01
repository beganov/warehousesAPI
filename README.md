# Lamoda Warehouses API
Тестовое задание на позицию Junior Golang Developer

Для запуска проекта:  
```
make up
```
Команда запускает БД PostgreSQL и приложение в Docker-контейнерах, применяет миграции и наполняет БД тестовыми данными.  
Для доступа к Swagger документации нужно перейти по ссылке:  
[SwaggerDocs URL](http://localhost:8080/swagger/index.html)

### Запросы к API:
1. Создание товара
```
curl -X 'PUT' \
  'http://localhost:8080/v1/items' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "items": [
    {
      "name": "TEST_SIZE",
      "quantity": 5,
      "size": "TEST_SIZE",
      "unique_id": "TEST0001",
      "warehouse_id": 0
    }
  ]
}'
```
**Ответ**: 
```
{
  "status": 201,
  "status_message": "Created",
  "message": "items successfully created",
  "error": ""
}
```
2. Создание склада
```
curl -X 'POST' \
  'http://localhost:8080/v1/warehouses' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "warehouse": {
    "availability": true,
    "name": "TEST_warehouse"
  }
}'
```
**Ответ**:
```
{
  "status": 201,
  "status_message": "Created",
  "message": "warehouse created successfully",
  "error": ""
}
```
3. Создание резервации
```
curl -X 'POST' \
  'http://localhost:8080/v1/reserve' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "ids": [
    "A0001",
    "A0002"
  ]
}'
```
**Ответ**:
```
{
  "status": 201,
  "status_message": "Created",
  "message": "reservation successfully created",
  "error": ""
}
```
4. Создание резервации (информаиця о товаре не записана в БД)
```
curl -X 'POST' \
  'http://localhost:8080/v1/reserve' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "ids": [
    "THIS_ITEM_IS_NOT_PRESENTED"  
]
}'
```
**Ответ**:
```
{
  "status": 201,
  "status_message": "Created",
  "message": "reservation successfully created",
  "error": ""
}
```
5. Удаление резервации
```
curl -X 'DELETE' \
  'http://localhost:8080/v1/reserve' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "ids": [
    "A0001"
  ]
}'
```
**Ответ**:
```
{
  "status": 200,
  "status_message": "OK",
  "message": "reservation successfully cancelled",
  "error": ""
}
```
6. Получение количества товаров на складе по идентификационному номеру
```
curl -X 'GET' \
  'http://localhost:8080/v1/items/0/quantity' \
  -H 'accept: application/json'
```
**Ответ**:
```
{
  "status": 200,
  "status_message": "OK",
  "message": {
    "TEST0001": 10,
    "THIS_ITEM_IS_NOT_PRESENTED": -1
  },
  "error": ""
}
```
7. Запрос на резервацию к недопступному складу:
```
curl -X 'POST' \
  'http://localhost:8080/v1/reserve' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "ids": [
    "F0001"
  ]
}'
```
**Ответ**:
```
{
  "status": 403,
  "status_message": "Forbidden",
  "message": null,
  "error": "warehouse is unavailable"
}
```

