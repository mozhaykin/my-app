# PromQL (Prometheus)

```shell
# RPS
sum(rate(http_server_requests_total[1m]))

# RPS по ручкам
sum(rate(http_server_requests_total[1m])) by (method)

# RPS статус 200
sum(rate(http_server_requests_total{status = "200"}[1m])) by (status)

# RPS ошибок
sum(rate(http_server_requests_total{status != "200"}[1m])) by (status)

# Процент ошибок от общего числа запросов
(
  sum(rate(http_server_requests_total{status !="200"}[1m]))
  /
  sum(rate(http_server_requests_total[1m]))
) * 100

# Duration среднее
sum(rate(http_server_request_duration_seconds_sum{}[1m]) / rate(http_server_request_duration_seconds_count{}[1m])) by (method)

# Duration 50-й перцентиль (медиана)
histogram_quantile(0.50, sum(rate(http_server_request_duration_seconds_bucket[1m])) by (le, method))

# Duration 90-й перцентиль
histogram_quantile(0.90, sum(rate(http_server_request_duration_seconds_bucket[1m])) by (le, method))

# Duration 95-й перцентиль
histogram_quantile(0.95, sum(rate(http_server_request_duration_seconds_bucket[1m])) by (le, method))

# Duration 99-й перцентиль
histogram_quantile(0.99, sum(rate(http_server_request_duration_seconds_bucket[1m])) by (le, method))

# Duration 99.9-й перцентиль
histogram_quantile(0.999, sum(rate(http_server_request_duration_seconds_bucket[1m])) by (le, method))
```

# TraceQL (Grafana Tempo)
Документация: https://grafana.com/docs/tempo/latest/traceql/

```shell
# Найти все трейсы дольше 5 сек
{duration > 5s}

# По имени трейса
{name="http GET /get"}

# По типу
{kind=server}

# По статусу OTEL
{status=unset}

# Resource: сквозные атрибуты для всех трейсов сервиса
{resource.service.name="server"}

# Event: имя
{event:name="error"}

# Event: атрибут
{event.error.message="Bad Request"}

# Span: атрибут
{span.http.request.method="GET"}

# Логическое И
{span.http.request.method="GET" && span.http.response.status_code=200}
{span.http.request.method="GET"} && {span.http.response.status_code=200}


# Логическое ИЛИ
{span.http.request.method="GET" || span.http.request.method="POST"}
{span.http.request.method="GET"} || {span.http.request.method="POST"}

# Если большой запрос можно так форматировать
{
  resource.service.name="server" &&
  name="http GET /get" &&
  span.http.response.status_code>=400
}

# Функции

# RPS
{resource.service.name="server"} | rate() by (span.http.response.status_code)

# RPS ошибок
{span.http.response.status_code>=400} | rate() by (span.http.response.status_code)

# Аналог increase
{resource.service.name="server"} | count_over_time() by (span.http.response.status_code)

# Квантиль
{span.http.response.status_code=200} | quantile_over_time(duration, .999, .99, .9)

# Гистрограмма
{span.http.response.status_code=200} | histogram_over_time(duration)

```
# LogQL (Grafana Loki)
Документация: https://grafana.com/docs/loki/latest/query/

```shell
# Базовый запрос
{compose_service="app", level="info"}

# = точное совпадение
# != не равно
# =~ совпадение с регулярным выражением
# !~ не совпадает с регулярным выражением

# Фильтрация строк по содержанию подстроки
{compose_service="app", level="error"} |= "server"

# |= содержит строку
# != не содержит строку
# |~ совпадает с regex
# !~ не совпадает с regex

# Работа с JSON
{app_name="my-app"} | json

# Фильтрация по лейблам и строкам
{app_name="my-app", level="error"} | json | method = "GET" and method != "POST"
# поддерживаемые операторы: and, or. У лейблов: =, !=, >, >=, <, <=

# RPS:
rate({app_name="my-app"}[1m])
sum(rate({app_name="my-app"}[1m])) by (level)
sum(rate({app_name="my-app"} | json | code != "200" [1m])) by (method, code)
sum(rate({app_name="my-app"} | json | level = "error" [1m])) by (method)
sum(rate({app_name="my-app"} | json | level = "error" |= "server" [1m])) by (method)
sum(rate({app_name="my-app"}[1m])) by (code)

# Квантиль
quantile_over_time(0.9,
  {app_name="my-app"}
    | json
    | code=~"400|500"
    | unwrap duration[1m]) by (method)

# Соотношение ошибок к общему числу запросов
sum(rate({app_name="my-app", level="error"}[1m])) / sum(rate({app_name="my-app"}[1m]))
```
