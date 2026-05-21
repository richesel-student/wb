# Утилита мониторинга GC

Go-утилита для анализа производительности сборщика мусора (GC) и памяти рантайма с экспортом данных в формате Prometheus.

## Возможности

- **Сбор метрик**: Периодически вызывает `runtime.ReadMemStats` и обновляет метрики Prometheus
- **Интеграция с Prometheus**: Отдаёт метрики через endpoint `/metrics`
- **Управление GC**: GET/POST endpoint `/gc-percent` для чтения или изменения `debug.SetGCPercent()` во время выполнения
- **Профилирование**: Интеграция с `net/http/pprof` для CPU, heap и goroutine профилей через `/debug/pprof/`

## Метрики

Утилита экспортирует следующие метрики Prometheus:

- `gc_monitor_mem_alloc_bytes` — Память, используемая в данный момент (`Alloc`)
- `gc_monitor_mem_total_alloc_bytes_total` — Общее количество выделенных байт (`TotalAlloc`)
- `gc_monitor_mem_sys_bytes` — Память, запрошенная у ОС (`Sys`)
- `gc_monitor_gc_cycles_total` — Количество циклов сборки мусора (`NumGC`)
- `gc_monitor_gc_last_time_seconds` — Время последнего GC в Unix timestamp (`LastGC`)
- `gc_monitor_gc_pause_total_ns` — Общее время пауз GC в наносекундах (`PauseTotalNs`)
- `gc_monitor_goroutines_count` — Количество активных goroutine (`runtime.NumGoroutine`)
- `gc_monitor_heap_objects_count` — Количество объектов в куче (`HeapObjects`)

## Быстрый старт

### Требования

- Go 1.21 или новее




### Установка

1. Клонируйте репозиторий
2. Перейдите в директорию проекта:

3. Загрузите зависимости:
   ```bash
   go mod tidy
   ```

### Запуск приложения

**С помощью `make`:**
```bash
make run
```

**С помощью `go run`:**
```bash
go run main.go
```

**Сборка и запуск бинарного файла:**
```bash
make build
./gc-monitor
```

Сервер будет запущен на порту 8080.

## Использование

### Просмотр метрик

Чтобы посмотреть метрики Prometheus:
```bash
curl http://localhost:8080/metrics
```

Пример вывода:
```
# TYPE gc_monitor_mem_alloc_bytes gauge
gc_monitor_mem_alloc_bytes 1234567
# TYPE gc_monitor_mem_total_alloc_bytes_total counter
gc_monitor_mem_total_alloc_bytes_total 4567890123
# TYPE gc_monitor_mem_sys_bytes gauge
gc_monitor_mem_sys_bytes 9876543210
# TYPE gc_monitor_gc_cycles_total counter
gc_monitor_gc_cycles_total 42
# TYPE gc_monitor_gc_last_time_seconds gauge
gc_monitor_gc_last_time_seconds 1680000000
# TYPE gc_monitor_gc_pause_total_ns counter
gc_monitor_gc_pause_total_ns 123456789
# TYPE gc_monitor_goroutines_count gauge
gc_monitor_goroutines_count 10
# TYPE gc_monitor_heap_objects_count gauge
gc_monitor_heap_objects_count 723456
```

### Управление GC Percent

**Получить текущее значение GC percent:**

```bash
curl http://localhost:8080/gc-percent
```

**Установить GC percent (например, 50):**
```bash
curl -X POST http://localhost:8080/gc-percent?value=50
```

GC percent определяет, как часто будет запускаться сборщик мусора.
Меньшие значения делают GC более агрессивным, большие — позволяют выделять больше памяти перед запуском GC.
### Использование pprof

Приложение включает `net/http/pprof` по умолчанию. Профили доступны по пути `/debug/pprof/`.

**Собрать heap profile:**
```bash
go tool pprof http://localhost:8080/debug/pprof/heap
```

**CСобрать CPU profile за 30 секунд:**
```bash
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
```

**Список всех доступных профилей**
```bash
curl http://localhost:8080/debug/pprof/
```

## Структура проекта

```
./L4.4/
├── go.mod            - Go module definition
├── go.sum            - Dependency checksums
├── main.go           - Application entry point and HTTP server
├── internal/
│   ├── metrics/
│   │   └── collector.go - Prometheus metrics collector
│   ├── handlers/
│   │   └── handlers.go  - HTTP handlers for GC percent control
├── Makefile          - Build and run commands
├── README.md         - Documentation
```



### Testing
```bash
make test
```

### Сборка
```bash
make build
```

### Очистка
```bash
make clean
```

