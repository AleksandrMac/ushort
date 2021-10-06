# ushort
Сервис сокращения ссылок.

## Heroku
https://mac-short.herokuapp.com/

## SwaggerHub
https://app.swaggerhub.com/apis/AleksandrMac/UShort/0.0.1

## Example
### Регистрация на сервисе 
```
curl -X POST https://mac-short.herokuapp.com/auth/sign-up -d '{"email":"first@user.ru", "password":"12345"}'
```

### Авторизация на сервисе 
```
curl -X POST https://mac-short.herokuapp.com/auth/sign-in -d '{"email":"first@user.ru", "password":"12345"}'
```
в ответ получите токен авторизации 
```
{"Authorization": "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTgwYzg0YTUtMTliMi00MmExLTk4NzUtNzE1YzBkNWNlYjRmIn0.4zWJ8puffcDwBXGDaiKVtIKWiSeaCmF8nsScA_VF_Sk"}
```
### Редирект 
```
curl -X GET https://mac-short.herokuapp.com/{urlID}
```

### Генерация ссылки
```
curl -X GET https://mac-short.herokuapp.com/url/generate -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOWMwZWY1NzEtMjg3Mi00MzViLWFhYzktZjBmNDAyZTZjYzliIn0.fvylVBxU8zYXth4dRwkFIdj6F0sckXRB11XentwBras"
```

### Получение списка ссылок
```
curl -X GET https://mac-short.herokuapp.com/url/ -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTgwYzg0YTUtMTliMi00MmExLTk4NzUtNzE1YzBkNWNlYjRmIn0.4zWJ8puffcDwBXGDaiKVtIKWiSeaCmF8nsScA_VF_Sk" -v
```

### Создание новой короткой ссылки
```
curl -X POST https://mac-short.herokuapp.com/url -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTgwYzg0YTUtMTliMi00MmExLTk4NzUtNzE1YzBkNWNlYjRmIn0.4zWJ8puffcDwBXGDaiKVtIKWiSeaCmF8nsScA_VF_Sk" -d '{"urlID": "besturl","redirectTo": "https://translate.yandex.ru","description": "instagram promo"}'
```

### Обновление информации о короткой ссылке
```
curl -X PATCH https://mac-short.herokuapp.com/url -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTgwYzg0YTUtMTliMi00MmExLTk4NzUtNzE1YzBkNWNlYjRmIn0.4zWJ8puffcDwBXGDaiKVtIKWiSeaCmF8nsScA_VF_Sk" -d '{"urlID":"besturl","redirectTo":"https://www.instagram.com/","description":"instagram promo"}'
```


### Получение короткой ссылки по ID
```
curl -X GET https://mac-short.herokuapp.com/url/besturl -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTgwYzg0YTUtMTliMi00MmExLTk4NzUtNzE1YzBkNWNlYjRmIn0.4zWJ8puffcDwBXGDaiKVtIKWiSeaCmF8nsScA_VF_Sk"
```

### Удаление короткой ссылки
```
curl -X DELETE https://mac-short.herokuapp.com/url/6qiri86cmq -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTgwYzg0YTUtMTliMi00MmExLTk4NzUtNzE1YzBkNWNlYjRmIn0.4zWJ8puffcDwBXGDaiKVtIKWiSeaCmF8nsScA_VF_Sk"
```





## Конфигурация
TMPURL_LIFE_TIME - время бронирования сгенерированой ссылки, по умолчанию 60 сек  
LENGTH_URL - длина сгенерированной ссылки  
LOG_LEVEL - уровень логирования, [trace, debug, info, warn, error, critical]  
SERVER_GRACEFULLTIME - время Graceful Shutdown  
DATABASE_URL - урл подключения к бд, по умолчанию: "postgres://postgres:password@localhost:5432/ushort?sslmode=disable"  
PORT - порт приложения, по умолчанию 8000  
SERVER_TIMEOUT - серверный тайм-аут  