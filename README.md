# service_sync_tables

## Сервис для синхронизации таблиц нескольких баз данных

### Переменные окружения необходимые для запуска сервиса:
 - URL_STORAGE_PASS нужен для указания пути к волуму внутри контейнера (имеет значение по умолчанию)
 - MONGO_HOST
 - MONGO_PORT
 - MONGO_COLLECTION
 - MONGO_DATABASE
 - API_SERVER_PORT
 - LOG_LEVEL (имеет занчение по умолчанию)

БД донор должен иметь элиас "main_db"