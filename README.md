# Run In K8s

Посмотреть на ArgoCD
```bash
http://argocd.goscl.ru
Username: admin
Password: F2OJameyuLoJca5T
```

В https://gitlab.golang-school.ru/potok-1 создайте папку, через кнопку `New subgroup` и назовите её своим юзернеймом.
Например: `mnepryakhin`

Внутри этой папки вы будете создавать свои приложения, например `my-app`

```bash
# Build
docker build -t harbor.goscl.ru/amozhaykin/my-app:v0.1.0 .

# Если у вас Мак, к билду добавьте флаг --platform="linux/amd64":
docker build -t harbor.goscl.ru/amozhaykin/my-app:v0.1.0 . --platform="linux/amd64"

# Push
docker push harbor.goscl.ru/amozhaykin/my-app:v0.1.0

# Если при пуше ошибка авторизации, залогинтесь в Harbor:
docker login -u "admin" -p "ahSpzTYk6Awxh7X" harbor.goscl.ru

# Убедитесь что образ загрузился в Harbor, зайдите в свой проект:
http://harbor.goscl.ru/harbor/projects

# Дальше опционально.
# Если хотите, можете создать своё приложение в ArgoCD и загрузить его в k8s.

# Нужно создать ArgoCD Application и манифесты k8s
# Клонируем к себе репозиторий
https://gitlab.golang-school.ru/potok-1/deploy

# Создайте свой файл как mnepryakhin-my-app.yaml (замените mnepryakhin на свой логин)
# Замените все mnepryakhin в файле на свое имя пользователя

# Дабавьте манифесты как в папке mnepryakhin-my-app (во всех файлах замените mnepryakhin на свой логин)
```

Далее пушите всё в gitlab.

В http://argocd.goscl.ru в течение 3 минут должен появится ваш проект.

Ваше приложение будет доступно по пути amozhaykin:
```bash
http://k8s.goscl.ru/mnepryakhin/my-app/hello
```