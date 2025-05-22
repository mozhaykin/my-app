# Запуск приложения локально
http://localhost:8080/amozhaykin/my-app/api/v1/profile

# Запуск приложения на сервере
http://k8s.goscl.ru/amozhaykin/my-app/api/v1/profile

# GET - получить профиль
принимает ID в URL запроса
возвращает JSON: {"id": ID, "name": "Alice", "age": 30} и код ответа 200

# POST - создать профиль
принимает JSON: {"name": "Bob", "age": 33}
возвращает JSON: {"id": ID} и код ответа 200

# PUT - обновить существующий или создать новый профиль
принимает JSON: {"id": ID, "name": "Alice", "age": 30}
возвращает код ответа 204

# DELETE - удалить профиль
принимает ID в URL запроса
возвращает код ответа 204
