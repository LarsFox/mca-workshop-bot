# Эксперт по оскорблениям

Телеграм-бот для Математического Воркшопа НГУ.

Для первого запуска используем `make init` и заполняем `.env` файл.

## Запускаем на локальной машине

Для имитации работы модели подставляем `http://localhost:11112` во все переменные `MCA_WORKSHOP_MODEL_*` и запускаем `make mock` в отдельной сессии.

## Запускаем в Докере

1. Собираем образ через сервера: `make docker-build`.
1. Логинимся в Гитлабе: `docker login registry.gitlab.com`.
1. Проверяем, что в `.env` указаны верные пути до контейнеров и запускаем `make docker-run`.

## Возможные проблемы

Если контейнеры с моделями неожиданно прекращают работу, попробуйте увеличить доступную память Докера.
