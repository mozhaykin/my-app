# Запуск приложения локально
```bash
Запрос GET http://localhost:8080/amozhaykin/my-app/api/v1/profile/ID
возвращает JSON: {"name": "Alice", "age": 30} и код ответа 200

Запрос POST http://localhost:8080/amozhaykin/my-app/api/v1/profile
принимает JSON: {"name": "Bob", "age": 33}, а возвращает код ответа 201
```

# Запуск приложения на сервере
```bash
Запрос GET http://k8s.goscl.ru/amozhaykin/my-app/api/v1/profile/ID
возвращает JSON: {"name": "Alice", "age": 30} и код ответа 200

Запрос POST http://k8s.goscl.ru/amozhaykin/my-app/api/v1/profile
принимает JSON: {"name": "Bob", "age": 33}, а возвращает код ответа 201
```