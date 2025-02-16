[![Go](https://img.shields.io/badge/-Go-464646?style=flat-square&logo=Go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-464646?style=flat-square&logo=PostgreSQL)](https://www.postgresql.org/)
[![docker](https://img.shields.io/badge/-Docker-464646?style=flat-square&logo=docker)](https://www.docker.com/)
[![HTTP](https://img.shields.io/badge/-HTTP-464646?style=flat-square&logo=http)](https://developer.mozilla.org/en-US/docs/Web/HTTP)

# merch_store
# Магазин мерча

---
## Описание проекта
Проект представляет собой магазин мерча, который позволяет сотрудникам обмениваться монетками и приобретать мерч за эти монетки. Каждый сотрудник может видеть список купленных товаров и историю перемещения монеток.

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
### Регистрация/Авторизация и получение токена:

**Request:**
```
curl -X POST http://localhost:8080/api/auth \
  -H "Content-Type: application/json" \
  -d '{"username": "<ваше_имя>", "password": "<ваш_пароль>"}'
```
**Response:**
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3Mzk3ODQzMzgsImlhdCI6MTczOTY5NzkzOH0.O_hKMSeORJxUY2e8sHSoeYMH6kN8Lfwc7nLAVgpA4zE"
}
```

### Получение информации о балансе, покупках и транзакциях:

**Request:**
```
curl http://localhost:8080/api/info \
  -H "Authorization: Bearer "<ваш_токен>""
```
**Response:**
```
{
    "coins": 700,
    "inventory": [
        {
            "type": "wallet",
            "quantity": 1
        },
        {
            "type": "book",
            "quantity": 1
        }
    ],
    "coinHistory": {
        "received": [],
        "sent": [
            {
                "toUser": "sergey",
                "amount": 100
            },
            {
                "toUser": "sergey",
                "amount": 100
            }
        ]
    }
}
```

### Отправка монет другому пользователю:

**Request:**
```
curl -X POST http://localhost:8080/api/sendCoin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer "<ваш_токен>"" \
  -d '{"toUser": "<имя_получателя>", "amount": <количество_монет>}'
```
**Response:**
```
{
    "toUser": "sergey",
    "amount": 100
}
```

### Проверим баланс после отправки монет:
```
{
    "coins": 650,
    "inventory": [
        {
            "type": "wallet",
            "quantity": 1
        },
        {
            "type": "book",
            "quantity": 1
        }
    ],
    "coinHistory": {
        "received": [],
        "sent": [
            {
                "toUser": "sergey",
                "amount": 50
            },
            {
                "toUser": "sergey",
                "amount": 100
            },
            {
                "toUser": "sergey",
                "amount": 100
            }
        ]
    }
}
```
```
{
    "coins": 950,
    "inventory": [
        {
            "type": "hoody",
            "quantity": 1
        }
    ],
    "coinHistory": {
        "received": [
            {
                "fromUser": "aleksey",
                "amount": 50
            },
            {
                "fromUser": "aleksey",
                "amount": 100
            },
            {
                "fromUser": "aleksey",
                "amount": 100
            }
        ],
        "sent": []
    }
}
```

### Покупка мерча:

**Request:**
```
curl http://localhost:8080/api/buy/<предмет_мерч> \
  -H "Authorization: Bearer "<ваш_токен>""
```
**Проверка купленных предметов:**
```
{
    "coins": 500,
    "inventory": [
        {
            "type": "pink-hoody",
            "quantity": 1
        }
    ],
    "coinHistory": {
        "received": [],
        "sent": []
    }
}
```

### Попытка купить несуществующий товар:

**Request:**
```
curl http://localhost:8080/api/buy/milk \
  -H "Authorization: Bearer "<ваш_токен>"
```
**Response:**
```
{
    "errors": "item not found"
}
```

### Попытка авторизации с пустым именем/паролем:

**Request:**
```
curl -X POST http://localhost:8080/api/auth \
  -H "Content-Type: application/json" \
  -d '{"username": , "password": "1234"}'
```
**Response:**
```
{
    "errors": "invalid request format"
}
```

## Нагрузочное тестирование
**Для нагрузочного тестирования использовал утилиту [plow](https://github.com/six-ddc/plow)**

### Нагрузочный тест на /api/info
```
plow http://localhost:8080/api/info -c 8 -H "Authorization: Bearer <ваш_токен>"
Benchmarking http://localhost:8080/api/info using 8 connection(s).
@ Real-time charts is listening on http://[::]:18888

Summary:
  Elapsed      1m22s
  Count        26686
    2xx        26686
  RPS        325.055
  Reads    0.104MB/s
  Writes   0.073MB/s

Statistics    Min     Mean     StdDev      Max
  Latency   2.994ms  24.59ms  17.502ms  168.797ms
  RPS       190.86   325.11     62.4     452.13

Latency Percentile:
  P50         P75       P90       P95       P99       P99.9     P99.99
  20.829ms  33.277ms  47.138ms  58.047ms  84.223ms  118.524ms  159.254ms

Latency Histogram:
  11.837ms   11784  44.16%
  29.371ms   11448  42.90%
  43.112ms    2312   8.66%
  60.967ms     742   2.78%
  82.982ms     301   1.13%
  101.54ms      78   0.29%
  119.905ms     13   0.05%
  140.121ms      8   0.03%
```
**Резюме теста (Summary):**
* **Elapsed: 1m22s** — общее время выполнения теста.
* **Count: 26686** — всего выполнено 26 686 HTTP-запросов.
* **2xx: 26686** — все ответы были с успешным статусом (2xx).
* **RPS (Requests per Second):** 325.055 — средняя скорость запросов была ~325 запросов в секунду.
* **Reads: 0.104MB/s** — скорость чтения данных (ответов) составляла ~0.1 МБ в секунду.
* **Writes: 0.073MB/s** — скорость записи данных (запросов) составляла ~0.07 МБ в секунду.

**Статистика (Statistics):**
**Latency (Задержка):**
* **Min: 2.994ms** — минимальная задержка.
* **Mean: 24.59ms** — средняя задержка.
* **StdDev: 17.502ms** — стандартное отклонение.
* **Max: 168.797ms** — максимальная задержка.

**RPS (запросов в секунду):**
* **Min: 190.86** — минимальная зафиксированная скорость запросов.
* **Mean: 325.11** — средняя скорость запросов.
* **StdDev: 62.4** — небольшие колебания скорости.
* **Max: 452.13** — максимальная скорость запросов.

![info](tests/load/info.png)

### Нагрузочный тест на /api/auth
```
plow http://localhost:8080/api/auth \
>   -c 8 \
>   -m POST \
>   -H "Content-Type: application/json" \
>   --body @auth_payload.json
Benchmarking http://localhost:8080/api/auth using 8 connection(s).

Summary:
  Elapsed     1m1.1s
  Count        38436
    2xx        38429
    5xx            7
  RPS        628.758
  Reads    0.159MB/s
  Writes   0.112MB/s

Statistics   Min      Mean     StdDev      Max
  Latency     0s    12.711ms  12.642ms  138.256ms
  RPS       295.93   628.4      132      831.98

Latency Percentile:
  P50        P75       P90       P95       P99      P99.9     P99.99
  6.978ms  20.479ms  30.301ms  37.004ms  53.842ms  85.638ms  118.076ms

Latency Histogram:
  5.189ms    21333  55.50%
  17.492ms   12390  32.24%
  31.182ms    3917  10.19%
  43.203ms     613   1.59%
  65.901ms     168   0.44%
  87.062ms      10   0.03%
  102.295ms      4   0.01%
  138.256ms      1   0.00%
```

![auth](tests/load/auth.png)

### Нагрузочный тест на /api/sendCoin
```
plow http://localhost:8080/api/sendCoin \
> -c 1 -n 1000 \
> -m POST \
> -H "Content-Type: application/json" \
> -H "Authorization: Bearer <ваш_токен>" \
> --body @send_coin_payload.json
Benchmarking http://localhost:8080/api/sendCoin with 1000 request(s) using 1 connection(s).
@ Real-time charts is listening on http://[::]:18888

Summary:
  Elapsed       3.9s
  Count         1000
    2xx         1000
  RPS        254.244
  Reads    0.018MB/s
  Writes   0.082MB/s

Statistics    Min     Mean    StdDev    Max
  Latency   2.998ms  3.924ms  683µs   19.149ms
  RPS       245.83   254.62    7.67    260.02

Latency Percentile:
  P50        P75      P90     P95      P99     P99.9     P99.99
  3.697ms  4.194ms  4.565ms  4.78ms  5.377ms  19.149ms  19.149ms

Latency Histogram:
  3.63ms   475  47.50%
  3.925ms  279  27.90%
  4.265ms  157  15.70%
  4.618ms   36   3.60%
  4.855ms   28   2.80%
  5.115ms    7   0.70%
  5.328ms   14   1.40%
  5.643ms    4   0.40%
```

![coin](tests/load/coin.png)

### Нагрузочный тест на /api/buy/socks
```
plow http://localhost:8080/api/buy/socks -c 1 -n 100 -H "Authorization: Bearer <ваш_токен>"
Benchmarking http://localhost:8080/api/buy/socks with 100 request(s) using 1 connection(s).
@ Real-time charts is listening on http://[::]:18888

Summary:
  Elapsed      400ms
  Count          100
    2xx          100
  RPS        243.927
  Reads    0.017MB/s
  Writes   0.056MB/s

Statistics    Min     Mean    StdDev     Max
  Latency   3.111ms  4.094ms  1.774ms  20.751ms

Latency Percentile:
  P50       P75      P90      P95      P99     P99.9     P99.99
  3.73ms  4.181ms  4.722ms  5.299ms  6.071ms  20.751ms  20.751ms

Latency Histogram:
  3.933ms  72  72.00%
  4.164ms   9   9.00%
  4.191ms   4   4.00%
  4.232ms   5   5.00%
  4.699ms   5   5.00%
  5.291ms   3   3.00%
  5.35ms    1   1.00%
  6.071ms   1   1.00%
```

![buy](tests/load/buy.png)

---
## Разработал:
[Aleksey Kazikov](https://github.com/KazikovAP)

---
## Лицензия:
[MIT](https://opensource.org/licenses/MIT)
