# Запуск приложения локально
http://localhost:8080/amozhaykin/my-app/api/v1/profile
http://localhost:8080/amozhaykin/my-app/api/v2/profile

# Запуск приложения на сервере
http://k8s.goscl.ru/amozhaykin/my-app/api/v1/profile
http://k8s.goscl.ru/amozhaykin/my-app/api/v2/profile

# POST - создать профиль
принимает JSON в теле запроса: {"name": "Bob", "age": 33}
возвращает JSON: {"id": ID} и код ответа 201

# GET - получить профиль
принимает ID в URL запроса
возвращает JSON: {"id": ID, "name": "Alice", "age": 30} и код ответа 200

# PUT - обновить существующий или создать новый профиль
принимает JSON в теле запроса: {"id": ID, "name": "Alice", "age": 30}
возвращает код ответа 204

# DELETE - удалить профиль
принимает ID в URL запроса
возвращает код ответа 204