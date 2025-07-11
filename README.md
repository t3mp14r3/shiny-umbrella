# shiny-umbrella

## Описание
Проект состоит из стандартного набора небольших модулей для обеспечения работы. API состоит из [handler](https://github.com/t3mp14r3/shiny-umbrella/tree/master/internal/api/handler)'a (1 уровень веб сервера) и [usecase](https://github.com/t3mp14r3/shiny-umbrella/tree/master/internal/api/usecase)'a (бизнес логика). Код взаимодействия с базой данных находится в [repository](https://github.com/t3mp14r3/shiny-umbrella/tree/master/internal/repository), там же лежит и [notifier](https://github.com/t3mp14r3/shiny-umbrella/tree/master/internal/repository/notifier), это модуль, который отвечает за принятие сообщений от Postgres через систему NOTIFY/LISTEN, нужно это для того, что отлавливать изменения в турнирах и вовремя обновлять их данные, а также (это случилось в процессе разработки), для мониторинга создания новых автоматических турниров в процессе работы приложения. Управляются турниры через модуль [cron](https://github.com/t3mp14r3/shiny-umbrella/tree/master/internal/cron), он производит первоначальное планирование турниров при запуске приложения, также, как и в последствии, во время работы приложения. Вот очень упрощенная схема работы турниров:
<img width="1256" height="660" alt="image" src="https://github.com/user-attachments/assets/fb912114-6ece-462e-97fa-1424747e653d" />


## Алиби
Большая часть кода появилась "стихийно", пробовал разные методы, комбинировал их получилось вот что-то такое, не могу сказать, что горжусь этой реализацией и само собой, этот код бы не пошел на прод под серьезную нагрузку, нужно использовать более серьезный планировщик и добавлять redundancy. Да и по большей своей части многие методы (вроде синхронизации cron'a и postgres'a) такие себе, можно попробовать смотреть на них, как на "креативные решения".

## Запуск
Для простоты, решил просто взять docker-compose, вся конфигурация находится в .env файле
```
docker-compose up --build
```

## Создание турниров
В конце-концов, решил сделать создание турниров не через API (с радостью объясню почему), создание происходит непосредственно в бд.

**Ручные турниры**<br>
```
BEGIN;
INSERT INTO tournaments(price, min_users, max_users, bets, starts_at, duration) VALUES(99, 2, 10, 5, NOW() + INTERVAL '30 minutes', INTERVAL '10 minutes');
INSERT INTO rewards(tournament_id, place, prize) VALUES(1, 1, 300),(1, 2, 200),(1, 3, 100);
COMMIT;
```

**Автоматические турниры**<br>
```
BEGIN;
INSERT INTO automatic(price, min_users, max_users, bets, duration, repeat) VALUES(99, 2, 10, 5, INTERVAL '30 minutes', INTERVAL '5 hours');
INSERT INTO automatic_rewards(automatic_id, place, prize) VALUES(1, 1, 300),(1, 2, 200),(1, 3, 100);
COMMIT;
```

## API
**Создание пользователя**<br>
`POST` `/users`
```
{
  "username": "JohnDoe",
  "balance": 100
}
```

**Получение пользователей**<br>
`GET` `/users`

**Редактирование пользователя**<br>
`PUT` `/users`
```
{
  "username": "JohnDoe",
  "balance": 250
}
```

**Получение турниров**<br>
`GET` `/tournaments`

**Получение турниров пользователя**<br>
`GET` `/tournaments?username=JohnDoe`

**Регистрация на турнир**<br>
`POST` `/tournaments/register`
```
{
  "username": "JohnDoe",
  "tournament_id": 1
}
```

**Отправка score'a на турнир**<br>
`POST` `/tournaments/score`
```
{
  "username": "JohnDoe",
  "tournament_id": 1,
  "score": 25
}
```
