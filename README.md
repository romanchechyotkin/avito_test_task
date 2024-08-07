# ESTATE MANAGEMENT SERVICE

[![Build Status][ci-badge]][ci-runs]
![coverage][coverage-badge]

[ci-badge]:            https://github.com/romanchechyotkin/avito_test_task/actions/workflows/test.yaml/badge.svg
[ci-runs]:             https://github.com/romanchechyotkin/avito_test_task/actions
[coverage-badge]:      https://raw.githubusercontent.com/romanchechyotkin/avito_test_task/badges/.badges/master/coverage.svg

Микросервис, с помощью которого пользователь сможет продать квартиру, загрузив объявление на Авито. Перед появлением объявления на сайте, квартира должна пройти модерацию

[Полное тестовое задание](https://github.com/avito-tech/backend-bootcamp-assignment-2024?tab=readme-ov-file)


Используемые технологии:
- PostgreSQL (база данных)
- Docker (для запуска сервиса)
- Swagger (для документации API)
- Gin (веб фреймворк)
- golang-migrate/migrate (для миграций БД)
- pgx (драйвер для работы с PostgreSQL)
- uber/mock, testify (для тестирования)

# Configuration

Сконфигурировать приложение можно используя `config.yaml`, указав путь до файла в переменной `CONFIG_PATH`
#### example
```shell
  export CONFIG_PATH=config/config.yaml
```

```yaml
http:
  host: localhost
  port: 8080
postgresql:
  user: user
  password: password
  host: localhost
  port: 5432
  database: postgres
  ssl_mode: disable
  auto_create: false
jwt:
  sign_key: secret
  token_ttl: 60m
```

вместе с файлом приложение можно настроить используя перемнные окружения
- `CONFIG_PATH=path` - настройка расположения yaml конфиг файла; дефолт значение = "config.yaml"
- `APP_ENV=prod/dev` - настройка окружения приложения
- `LOG_LEVEL=debug/info/warn/error` - настройка уровня логирования; дефолт значение = "debug"

для всех полей из yaml файла есть переменные окружения для конфигурации

- `PORT`
- `HOST`
- `PG_USER`
- `PG_PASSWORD`
- `PG_HOST`
- `PG_PORT`
- `PG_DATABASE`
- `PG_SSL`
- `PG_AUTO_CREATE`
- `JWT_KEY`
- `JWT_TTL`

# Usage

Запустить сервис можно с помощью команды 

```shell
   make compose-up
```

Документацию после завпуска сервиса можно посмотреть по адресу `http://localhost:8080/swagger/index.html`
с портом 8080 по умолчанию

Запуск юнит тестов
```shell
  make unit-test
```

Запуск интеграционных тестов
```shell
  make integration-test
```

Запуск всех тестов
```shell
  make test
```

Запуск всех тестов с покрытием для получения отчёта в html
```shell
  make coverage-html
```

# Decisions <a name="decisions"></a>

В ходе разработки был сомнения по тем или иным вопросам, которые были решены следующим образом:

1. Какой тип данных использовать для primary key в таблицах `houses` и `flats` в базе данных  
> Изначально выбор стоял между UUID и числовым типом. Известно, что UUID плохо влияет на перфоманс БД. Так как у нас по условии  
> `Всего квартир (до 20kk)` я решил использовать Serial для сущностей Дом и Квартира, чтобы в рамках проекта не заострять внимание на генерации уникального id. 

2. Использования билдера для запросов в БД (`Squirrel`)
> Я предпочитаю писать Raw SQL запросы, поэтому не использовал билдер запросов. В рамках проекта было одно место, куда бы хорошо устроился билдер, а именно
>
> ```go
> q := "SELECT id, number, house_id, price, rooms_amount, moderation_status, created_at, updated_at FROM flats WHERE house_id = $1"
>
>	if userType == "client" {
>		q += "AND moderation_status = 'approved'"
>	}
> ```
> Посчитал, что не стоит переписывать все на Squirrel из-за единичного случая

3. Как отслеживать, что квартира на модерации?
> Я добавил в таблицу `flats` поле `modearator_id`, чтобы указывать на каком модераторе висит квартира, 
> но потом решил, что буду доставать из квариры ее статус, 
> что позволит нам отдавать более детализированные ошибки на клиент, например: 
> 
> что нельзя возвращать квартиры на модерацию из approved и declined статусов, 
> но из on moderation можно, мало ли модератор не смог завершить модерацию по уважительным причинам 

4. Способ написания проекта
> Изначально хотелось написать проект, используя uber.fx и DI-контейнеры, потому что именно такой подход я увидел на первой работе, 
> но такой подход максимально не Go way из-за непрозрачного прокидывания зависимостей. 
> Именно поэтому я решил, что буду использовать стандартный подход для написания Go приложений, где все зависимости прокидывваются явно 
