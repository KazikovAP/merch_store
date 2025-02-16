[![Go](https://img.shields.io/badge/-Go-464646?style=flat-square&logo=Go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-464646?style=flat-square&logo=PostgreSQL)](https://www.postgresql.org/)
[![docker](https://img.shields.io/badge/-Docker-464646?style=flat-square&logo=docker)](https://www.docker.com/)
[![HTTP](https://img.shields.io/badge/-HTTP-464646?style=flat-square&logo=http)](https://developer.mozilla.org/en-US/docs/Web/HTTP)

# merch_store
# Магазин мерча

---
## Описание проекта
Проект представляет собой 

---
## Компоненты системы
1. Backend-сервис:

---
## Технологии
* Go 1.23.0
* PostgreSQL
* Docker
* REST API

---
## Запуск проекта

**1. Клонировать репозиторий:**
```
git clone https://github.com/KazikovAP/merch_store.git
```

**2. Сборка и запуск проекта:**
```
docker-compose up --build
```

**3. Остановка и удаление контейнеров:**
```
docker-compose down
```

## Примеры запросов к API
чтобы получить токен:
```
curl -X POST http://localhost:8080/api/auth \
  -H "Content-Type: application/json" \
  -d '{"username": "ваше_имя", "password": "ваш_пароль"}'
```

После получения токена можно сделать запрос к /api/info так:
```
curl http://localhost:8080/api/info \
  -H "Authorization: Bearer <ваш_токен>"
```


для отправки монет:
```
curl -X POST http://localhost:8080/api/sendCoin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ваш_токен>" \
  -d '{"toUser": "имя_получателя", "amount": 10}'
```

А для покупки предмета, например, hat:
```
curl http://localhost:8080/api/buy/hat \
  -H "Authorization: Bearer <ваш_токен>"
```



---
## Разработал:
[Aleksey Kazikov](https://github.com/KazikovAP)

---
## Лицензия:
[MIT](https://opensource.org/licenses/MIT)
