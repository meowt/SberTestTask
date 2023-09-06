# SberTestTask

## Задание
Сервис представляет собой хранилище JSON-объектов с HTTP-интерфейсом. Сохраненные объекты размещаются в оперативной памяти, имеется возможность задать время жизни объекта. Необходимо обеспечить запись содержимого хранилища в файл на диске и восстановление состояния хранилища из файла при запуске приложения. 

## Реализация
В качестве хранилища ~~использован Redis, т.к. он является кэш-хранилищем (хранит данные в оперативной памяти, с возможностью указывать время жизни объекта). Также, имеет возможность сохранять локально. Redis запускается из docker-compose.yml~~ использована собственная структура данных customMap:
```go
// CMap - custom map with concurrency support & TTL mechanism
type CMap struct {
	M  map[string]Value
	Mu sync.Mutex
}

type Value struct {
	V         interface{} `json:"value"`
	ExpiresAt time.Time   `json:"expires_at"`
}
```

Структура поддерживает конкурентную запись и чтение.

Для неё реализовано 4 метода: 
1. ```Put(key string, v interface{}, ttl time.Duration)``` - добавляет в мапу пару ключ-значение; при нулевом ttl, запись не имеет ограничения по продолжительности жизни;
2. ```Get(key string) (value interface{})``` - возвращает значение по ключу; при отсутствии записи - интерфейс быть равен nil;
3. ```StoreToFile(f *os.File, interval time.Duration) (err error)``` - запускает тикер с указанным интервалом, который будет записывать текущее состоянии мапы в переданный .json файл;
4. ```LoadFromFile(f *os.File) (*CMap, error)``` - загружает состояние мапы из переданного .json файла.

Также, в пакете присутствует функция ```New(deleteTick time.Duration) *CMap``` - запускает тикер для очищения от истекших записей мапу и возвращает указатель на структуру CMap.

Для http сервера был выбран пакет gin. Для получения метрик - prometheus.

## API Endpoints
| Endpoint Name | HTTP Method | URL                             | Description              |
|---------------|-------------|---------------------------------|--------------------------|
| getObject     | GET         | localhost:9876/objects/{key}    | Get object from storage  |
| putObject     | PUT         | localhost:9876/objects/{key}    | Write object into storage<br>Optional header "Expires" with int value (in milliseconds) |
| liveness      | GET         | localhost:9876/probes/liveness  | Check liveness status    |
| readiness     | GET         | localhost:9876/probes/readiness | Check readiness status   |
| metrics       | GET         | localhost:9876/metrics          | Get metrics              |

