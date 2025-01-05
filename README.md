# README для Web Приложения *Task Scheduler*

## Описание проекта

*Task Scheduler* — это простое и интуитивно понятное веб-приложение для планирования задач, разработанное для помощи пользователям в организации их повседневных дел. С помощью нашего приложения вы сможете:

- создавать и редактировать задачи
- устанавливать цикл повтора и сроки выполнения
- изменять статус задач 
- организовывать задачи по дате выполнения
- управлять задачами с помощью удобного интерфейса

## Что реализовано в приложении?

- ✔️ Реализован функционал создания задач
- ✔️ Добавлена возможность редактирования задач
- ✔️ Интеграция баз данных для хранения задач
- ✔️ Простой и привлекательный интерфейс
- ✔️ Реализована функция поиска и удаления задач

## Инструкция по запуску кода локально
#### 1. Перед запуском необходимо настроить переменные окружения, находящиеся в корне проекта в файле .env:
- укажите желаемый порт прослушивания сервера в переменной `TODO_PORT`
- укажите свой пароль для авторизации в переменной `TODO_PASSWORD`
#### 2. Если вы находитесь в корне проекта выполните команду для запуска приложения `go run cmd/main.go`.

#### Теперь вы можете открыть приложение в браузере по адресу: [http://localhost:7540](http://localhost:7540/ "Порт указываете тот, который укзан в TODO_PORT"). 

## Инструкция по запуску тестов

#### Для запуска тестов выполните команду `go test ./tests/...`.

## Инструкция по сборке и запуску проекта через Docker

Если вы хотите запустить проект через Docker, следуйте этим шагам:

1. Убедитесь, что у вас установлен Docker.
2. Для сборки docker образа проекта необходимо находиться в корне проекта.
3. Соберем образ, выполнив эту команду в терминале:
```
docker build --tag scheduler-app:v1 .
```
4. Теперь, когда образ собран можно запускать контейнер на внешнем порту **7540** выполнив команду в терминале:
```
docker run -it -p 7540:7540 scheduler-app:v1
```

 После запуска вы можете получить доступ к приложению по адресу: http://localhost:7540/login.html. Пароль для авторизации указывается тот, который задан в переменной окружения `TODO_PASSWORD`.

# Заключение
#### Спасибо за использование *Task Scheduler*. Надеюсь, что приложение поможет вам более организованно управлять вашим временем и задачами! 