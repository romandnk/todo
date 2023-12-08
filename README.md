
## Подготовка к запуску

1. Создайте файл `.env` в папке `/config` с необходимыми данными. Пример содержания файла:

    ```env
    POSTGRES_HOST=postgres
    POSTGRES_PORT=5432
    POSTGRES_USER=test
    POSTGRES_PASSWORD=1234
    POSTGRES_DB=todo_db
    POSTGRES_SSLMODE=disable
    
    HTTP_SERVER_HOST=0.0.0.0
    HTTP_SERVER_PORT=8080
    ```

## Запуск

### Запуск тестов и приложения
```bash
make full-run
```
### Запуск приложения
```bash
make run
```
### Запуск тестов
```bash
make test
```
