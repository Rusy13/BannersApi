# Сервис баннеров

Сервис для работы с баннерами пользователей, позволяющий создавать и удалять баннеры по фиче или тегу, обновлять их содержимое, а также возвращаться к предыдущим версиям баннера


---
Содержимое
---
- [Стек технологий](#technology_stack)
- [Как запуститься](#getting_started)
- [Использование](#usage)
- [Примеры](#examples)
- [Решения](#decisions)
- [Дополнительные примечания](#additional_notes)


---
# Стек технологий <a name="technology_stack"></a>
* [![Gin](https://img.shields.io/badge/Golang-blue?style=plastic&logoColor=yellow&logo=golang_logo)](https://go.dev/)
* [![Postgres](https://img.shields.io/badge/PostgreSQL-4169E1?style=plastic&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
* [![Docker](https://img.shields.io/badge/Docker-white?style=plastic&logo=docker&logoColor=2496ED)](https://www.docker.com/)
* [![pgx](https://img.shields.io/badge/Driver_pgx-blue?style=plastic&logo=adminer&logoColor=white)](https://pkg.go.dev/github.com/jackc/pgx)

Для гибкости и удобства в тестировании и масштабировании проекта, была выбрана Clean архитектура.


# Как запуститься <a name="getting_started"></a>

Для запуска сервиса необходимо предварительно:
* Создать самоподписные сертификаты и поместить их в internal/app

# Использование <a name="usage"></a>

Запустить БД, кафку и memcached через docker-compose up --build

Само приложение запускается из директории cmd/app командой go run main.go


Для запуска тестов необходимо выполнить команду make integration-test


# Examples <a name="examples"></a>

Для неограниченного доступа к данным нужно использовать token admin_token в headers
Некоторые примеры запросов
* [Создание баннера:POST https://localhost:9000/banner c телом:]
```
{
  "tag_ids": [61, 71, 81, 1, 3],
  "feature_id": 2,
  "content": {
    "title": "New Banner",
    "text": "Some text",
    "url": "http://example.com"
  },
  "is_active": true
}
```
* [Получение баннера для пользователя: GET https://localhost:9000/user_banner?tag_id=1&feature_id=9&use_last_revision=true]
* [Обновление содержимого баннера: PATCH https://localhost:9000/banner/5 c телом:]
```
  {
  "tag_ids": [1, 2, 3],
  "feature_id": 10,
  "content": {
    "title": "Updated Title4",
    "text": "Updated Text4",
    "url": "http://updated-example.com4"
  },
  "is_active": true
}
```


# Решения <a name="decisions"></a>

# Дополнительные примечания <a name="additional_notes"></a>

