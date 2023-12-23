# Итоговый проект по курсу "Разработка микросервисов на Go"

### Выполнили

- Шубников Андрей
- Осташков Максим

### Задание

Подробное описание задания можно найти на [anonimpopov/final-project](https://github.com/anonimpopov/final-project)

### Сборка проекта

Проект собирается с помощью docker-compose.yaml файла, по-умолчанию он использует переменные окружения из .env файла (
для работы внутри докера), если хотим запустить приложение локально, то специально добавили файл .env.local, где указаны
правильные IP-адрса и другие настройки окружения
Не очень поняли, зачем нужен .env.dev файл (предположу, что это как раз .env для докера; для удоной кастомизации мы
добавили аргумент -env в оба сервиса, чтобы можно было явно задавать откуда парсить ключи; в docker-compose тоже можно
задать файл)

#### Docker:

```shell
docker-compose up -d
```

#### Local:

```shell
./service -env=.env.local
```

#### Вспомогательные сервисы

- MongoDB, port: 27017
- Postgres, port: 5432
- Prometheus, port: 9090
- Jaeger Tracing, port: 16686
- Grafana, port: 3000 (login – admin, password – grafana)

#### Наши сервисы

- driver, port 8081
- location, port 8080

### Информация

Все сервисы покрыты тестами, используется трассировка запросов, метрики и логирование. Все это можно посмотреть в
Grafana, Prometheus и Jaeger соответственно. Также в Grafana можно посмотреть на дашборды с метриками и трассировкой
запросов. Ошибки отправляются в Sentry, логи в stdout.

Весь проект разрабатывался вместе, поэтому вклад каждого в каждый сервис примерно одинаковый.

Планируемая оценка – 8 баллов. 