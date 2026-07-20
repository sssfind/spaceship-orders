# Spaceship Orders

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/sssfind/7e5bafd7b0307b4c541657a022f55c36/raw/coverage.json)
![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-blue)
![License](https://img.shields.io/badge/license-MIT-green)

**Spaceship Orders** - это распределённая микросервисная платформа для заказа, оплаты и сборки космических кораблей. Проект спроектирован по принципам чистой архитектуры

---

## Архитектура системы

Взаимодействие компонентов системы построено на сочетании синхронного (gRPC) и асинхронного (Apache Kafka) подходов:

```mermaid
graph TD
    User([Пользователь]) -->|HTTP REST| OrderService[Order Service]
    User -->|gRPC| IAMService[IAM Service]
    
    OrderService -->|gRPC| IAMService
    OrderService -->|gRPC| InventoryService[Inventory Service]
    OrderService -->|gRPC| PaymentService[Payment Service]
    
    OrderService -->|Publish order-paid| Kafka[(Kafka Broker)]
    Kafka -->|Consume order-paid| AssemblyService[Assembly Service]
    AssemblyService -->|Publish ship-assembled| Kafka
    Kafka -->|Consume ship-assembled| OrderService
    
    %% Databases
    OrderService -->|pgx/v5| PostgresOrder[(PostgreSQL Orders)]
    IAMService -->|lib/pq| PostgresIAM[(PostgreSQL IAM)]
    IAMService -->|redigo| RedisIAM[(Redis Sessions)]
    InventoryService -->|mongo-driver| MongoInventory[(MongoDB Inventory)]
    
    %% Styling
    classDef service fill:#1e1e2f,stroke:#3b3b4f,stroke-width:2px,color:#fff;
    classDef db fill:#2d1b4e,stroke:#4a327a,stroke-width:2px,color:#e2d9f3;
    classDef broker fill:#0f3d59,stroke:#176087,stroke-width:2px,color:#fff;
    classDef client fill:#1b4d3e,stroke:#2d725f,stroke-width:2px,color:#fff;
    
    class OrderService,InventoryService,IAMService,PaymentService,AssemblyService service;
    class PostgresOrder,PostgresIAM,RedisIAM,MongoInventory db;
    class Kafka broker;
    class User client;
```

---

## Технологический стек

Проект разработан с использованием современного стека технологий, ориентированного на производительность и строгую типизацию контрактов:

| Слой / Компонент         | Технология                              | Описание                                                                    |
| :----------------------- | :-------------------------------------- | :-------------------------------------------------------------------------- |
| **Язык разработки**      | **Go 1.26**                             | Использование новейших возможностей языка для высокой производительности.   |
| **Синхронный транспорт** | **gRPC & REST**                         | Внутреннее взаимодействие по gRPC, внешнее по HTTP                          |
| **Генерация кода API**   | **OpenAPI (Ogen) & Protobuf (Buf)**     | Схемы являются единственным источником правды                               |
| **Брокер сообщений**     | **Apache Kafka (Sarama)**               | Асинхронный event-driven пайплайн для сборки кораблей                       |
| **Базы данных**          | **PostgreSQL (pgx/v5), MongoDB, Redis** | Оптимальный выбор СУБД под задачи каждого отдельного микросервиса           |
| **Качество кода**        | **Golangci-lint, Gofumpt, Gci**         | Строгие правила форматирования кода и статического анализа                  |
| **Тестирование**         | **Testcontainers-go, Mockery**          | Компонентные тесты в изолированных контейнерах Docker и автогенерация моков |
| **Оркестрация**          | **Docker Compose, Taskfile**            | Удобное локальное окружение и автоматизация рутинных действий               |

---

## Структура репозитория

Проект организован в виде монорепозитория (Go Workspaces):

*   [assembly/](./assembly) — микросервис сборки кораблей (эмулирует физический процесс сборки, потребляя события из Kafka)
*   [iam/](./iam) — сервис идентификации и аутентификации (регистрация, логин, сессии в Redis)
*   [inventory/](./inventory) — каталог запчастей и комплектующих для космических кораблей (хранилище MongoDB)
*   [order/](./order) — ядро системы: создание заказов, расчет стоимости, API для клиентов
*   [payment/](./payment) — сервис-заглушка для имитации транзакций оплаты
*   [platform/](./platform) — общая разделяемая библиотека (логирование, Graceful Shutdown, обертки над gRPC/Kafka/Redis)
*   [shared/](./shared) — общие схемы Protobuf/OpenAPI и сгенерированный на их основе код
*   [deploy/](./deploy) — инфраструктурные манифесты, Dockerfile и конфигурационные файлы Docker Compose

---

##  Быстрый старт

### Требования
Для локального запуска вам понадобятся:
1. **Docker & Docker Compose**
2. **Go 1.26** (для локальной разработки)
3. **Task** (`go-task` для запуска автоматизированных сценариев)

###  Генерация переменных окружения

```bash
task env:generate
```

### Запуск микросервисов
```bash
task up-all
```

### Остановка всех контейнеров
```bash
task down-all
```

---

## Тестирование и проверка качества кода

В проекте настроен надежный пайплайн тестирования и статического анализа:

### Юнит-тестирование
Запуск быстрых изолированных юнит-тестов для всех модулей:
```bash
task test
```

### Расчёт покрытия (Code Coverage)
Запуск тестов с вычислением тестового покрытия бизнес-логики (`service` и `repository` слоёв) с выводом красивого HTML-отчета:
```bash
# Генерация консольного отчёта
task test-coverage

# Генерация HTML и автоматическое открытие в браузере
task coverage:html
```

### Проверка стиля и Линтинг
Автоматическое форматирование кода (сортировка импортов, приведение к единому стилю) и запуск линтеров:
```bash
# Форматирование gofumpt + gci
task format

# Статический анализ golangci-lint
task lint
```

### Сценарные тесты
Вы можете запустить готовые сценарии проверки авторизации и межсервисного взаимодействия:
```bash
# Тест сценария регистрации, входа и получения профиля в IAM
task test-iam

# Тест сквозного создания заказа с авторизацией сессии
task test-api-with-auth
```

---

## Кодогенерация

Если вы изменили контракты `.proto` в [shared/proto](./shared/proto) или OpenAPI схемы в [shared/api](./shared/api):

1. Установите необходимые утилиты (`buf`, `ogen`, `yq`, `mockery` и др.) и сгенерируйте код:
   ```bash
   task gen
   ```
2. Обновите моки для юнит-тестов:
   ```bash
   task mockery:gen
   ```
3. Выполните синхронизацию зависимостей во всех модулях:
   ```bash
   task deps:update
   ```


## Статус разработки

Проект разрабатывается итеративно. На данный момент полностью реализована микросервисная база и событийная модель взаимодействия (завершены первые 6 недель дорожной карты)

- [x] **Фаза 1–2: Базовые микросервисы и CI/CD**
   - Реализовано синхронное взаимодействие по HTTP/REST и gRPC (сервисы Order, Inventory, Payment)
   - Внедрена Clean Architecture, написаны интеграционные тесты с Testcontainers
   - Настроены пайплайны GitHub Actions для проверки качества кода
- [x] **Фаза 3: Event-Driven архитектура**
   - Поднят брокер Apache Kafka (KRaft)
   - Настроены консьюмеры и продюсеры для асинхронных процессов (Assembly Service, Notification Service)
   - Окружение полностью упаковано в Docker Compose
- [ ] **Фаза 4: Observability (В процессе)**
   - Инструментирование Go-кода: сбор бизнес-метрик (Prometheus), распределенная трассировка (OpenTelemetry + Jaeger) и централизованное логирование
- [ ] **Фаза 5: Infrastructure as Code (IaC)**
   - Автоматизация развертывания отказоустойчивого кластера БД (HA PostgreSQL + Patroni) в облаке с помощью Terraform и Ansible
- [ ] **Фаза 6: Kubernetes & GitOps**
   - Развертывание K8s-кластера, упаковка сервисов в Helm-чарты, деплой через ArgoCD и управление секретами через SOPS
- [ ] **Фаза 7: Chaos Engineering & SRE**
   - Проведение контролируемых сбоев с помощью Chaos Mesh и генерация нагрузок k6 для проверки показателей RTO/RPO
