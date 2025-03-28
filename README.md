```bash
Запрос GET http://localhost:8080/amozhaykin/my-app/hello
возвращает строку "Hello!"

Запрос GET http://localhost:8080/amozhaykin/my-app/profile
возвращает JSON: {"name": "Alice", "age": 30} и код ответа 200

Запрос POST http://localhost:8080/amozhaykin/my-app/profile
принимает JSON: {"name": "Bob", "age": 33}, а возвращает код ответа 201
```