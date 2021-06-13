# ushort
Сервис сокращения ссылок и подсчета переходов по ним. Для обработки csv используется https://github.com/AleksandrMac/csv_query

## /cmd/ushort
Каталог командных файлов
- /ushort.go - Инициализация сервиса
- /redir.go - Перенаправление по ссылкам
- /generate.go - Генерация коротких ссылок
- /statistics.go - Получение статистики по ссылкам
## /store
Каталог хранения данных 
- /url.csv - сопоставление коротких и настоящих ссылок. [Структура ссылки](https://www.bing.com/search?q=structure+url&qs=n&form=QBRE&sp=-1&pq=structur+url&sc=8-12&sk=&cvid=89E02D1A140744E7A56C3C79587A0D20) Столбцы:
    - short - короткая ссылка (path)
    - full - полная ссылка (protocol+domain+path)
    - status - open or close
    - created - дата создания ссылки
    - closed - дата закрытия ссылки
    - desc - описание ссылки
- /statistics.csv - исторя переходов по ссылкам. Столбцы:
    - url - короткая ссылка
    - datetime - дата/время перехода по ссылке
- /blacklist.csv - Список заблокированных сайтов. 
    - domain - домен подозрительного сайта
    - datetime - дата/время добавления в список
    - desc - причина добавления
- /warninglist.csv - список жалоб
    - domain - домен подозрительного сайта
    - full - страница с жалобой
    - desc - описание проблемы
    - datetime - дата/время добавления жалобы

# Run
Запуск производить из корневого каталога проекта

    go run cmd/ushort/main.go
