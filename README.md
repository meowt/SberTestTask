# SberTestTask

## Задание
Сервис представляет собой хранилище JSON-объектов с HTTP-интерфейсом. Сохраненные объекты размещаются в оперативной памяти, имеется возможность задать время жизни объекта. Необходимо обеспечить запись содержимого хранилища в файл на диске и восстановление состояния хранилища из файла при запуске приложения. 

## Реализация
В качестве хранилища использован Redis, т.к. он является кэш-хранилищем (хранит данные в оперативной памяти, с возможностью указывать время жизни объекта). Также, имеет возможность сохранять локально.

Redis запускается из docker-compose.yml

Для http сервера был выбран пакет gin. Для получения метрик - prometheus.

## API Endpoints
| Endpoint Name | HTTP Method | URL                           | Description              |
|---------------|-------------|------------------------------|--------------------------|
| getObject     | GET         | localhost:9876/objects/{key}   | Get object from storage  |
| putObject     | PUT         | localhost:9876/objects/{key}   | Write object into storage<br>Optional header "Expires" with int value (in milliseconds) |
| liveness      | GET         | localhost:9876/probes/liveness  | Check liveness status    |
| readiness     | GET         | localhost:9876/probes/readiness | Check readiness status   |
| metrics       | GET         | localhost:9876/metrics       | Get metrics              |

