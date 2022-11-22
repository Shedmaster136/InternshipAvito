# InternshipAvito
## Микросервис управления счетами и заказами пользователей

При выполнении задания ставил перед собой цель как можно более подробно разобраться в структуре и работе микросервисов и их их интерфейсами. По этой причине старался избегать использования фреймворков и различных библиотек-оберток. 

## Этапы Проектирования микросервиса:
1. Проектирование базы данных.
2. Проектирование интерфейсов.
3. Выстраивание архитектуры проекта.

### 1. Проектирование базы данных
SQL код базы данных содержиться в файле SQL DB User Wallet.
Схема ниже.
![Untitled](https://user-images.githubusercontent.com/96218277/202832114-82297652-47bb-4471-ad20-338e155be3b4.png)

Было выделено 4 сущности: 
1. Пользователь - содержит ID и сведения о балансе
2. Услуга - содержит сведения об оказываемой пользователю услуге
3. Заказ - содержит сведения о типе услуги и пользователе, реализует связь Пользователь - Услуга мн ко мн.
4. Статус заказа - согласто техническому заданию, заказ может иметь несколько этапов: оплата, утверждение, отмена. О дате наступления каждого этапа необходимо хранить информацию в базе данных для реализации бухгалтерского отчета.

### 2. Проектирование интерфейсов
При проектировании HTTP интерфейсов старался придерживаться архитектуры REST API.
+ HTTP Метод:URI:назначение
+ GET:/wallet/:Получить остаток на кошельке пользователя
+ Возвращаемые значения:
{
 + "uid": $идентификатор пользователя
 + "balance": $остаток на счету
}
Статус 500 - в случае если сведений о пользователе еще нет в базе данных
+ Обязательные заголовки:
+ Content-Type = application/json
+ Тело:
{
 + "uid": $идентификатор пользователя,
}

+ PUT:/wallet/:Пополнение кошелька пользователя
   + Возвращаемые значения:
Cтатус 201 - в случае успеха
Статус 500 - в случае если сведений о пользователе еще нет в базе данных
   + Обязательные заголовки:
Content-Type = application/json
    + Тело:
{
    +  "uid": $идентификатор пользователя,
    +  "balance"; $сумма поплнения
}

+ PUT:/transaction/:Создать заказ и зарезервировать деньги из кошелька пользователя
    + Возвращаемые значения:
Статус 201 - в случае успеха
Статус 500 - в случае, если пользователя нет в базе / у него недостаточно средств
    + Обязательные заголовки
Content-Type = application/json
    + Тело:
{
        + "uid": $идентификатор пользователя,
        + "sid": $идентификатор услуги,
        + "oid": $идентификатор заказа,
        + "price": $cумма, которая резервируется,
}

+ UPDATE:/transation/: подтвердить списание денег у пользователя
Статус 201 - в случае успеха
Статус 500 - в случае, если пользователя нет в базе / у него недостаточно средств
     + Content-Type = application/json
     + Тело:
{
        + "uid": $идентификатор пользователя,
        + "sid": $идентификатор услуги,
        + "oid": $идентификатор заказа,
}

+ DELETE:/transation/: отменить списание денег у пользователя
    + Content-Type = application/json
    + Тело:
{
        + "uid": $идентификатор пользователя,
        + "sid": $идентификатор услуги,
        + "oid": $идентификатор заказа,
}
Две последние функции не меняют заказ пользователя, они создают дополнительные экземпляры сущности "Статус заказа", в которой фиксируется статус и время получения этого статуса.
### 3. Выстраивание архитектуры проекта
Много усилий было потрачено на то, чтобы понять базовую архитектуру микросервиса на языке golang. Язык Golang новый для меня и его особенности реализации принципов ООП выбили меня из колеи. 

Сервис был разбит на 4 пакета:
+ wserver - Пакет, реализующий инициализацию обработчиков, запуск и остановку сервера
+ handlers - Пакет, реализующий обработки HTTP запросов, изьятие аргументов из формата json и сохранение в структуры.
+ businesslogic - Пакет, реализующий логику приложения
+ sqllayer - Пакет, реализующий доступ к базе данных

#### wserver
Данные для инициализации: настройки сервера и карта (map) со ссылками на обработчики http запросов.
Задача: Принимать запросы и передавать обработчикам
+ Основной класс:
+ Server 
    * Методы:
    * Start() - регистрирует переданные функции обработчиков с соответствующими маршрутами, запускает сервер с заданными параметрами.
    * Stop() - graceful завершение работы сервера (go потоки закрываются по мере завершения последних принятых вызовов).

#### handlers
Данные для инициализации: базовая структура пакета businesslogic для вызова методов уровня логики приложения
Задача: Проверять формат запросов на соответсвие спецификации http интерфейсов (см. выше)
+ Основной класс:
- VerbHandler - реализует интерфейс http.Handler, чтобы сервер мог зарегистрировать как обработчик события. содержит ссылку на базовый класс уровня логики
  * Методы:
  * ServeHTTP() - метод интерфейса
+ Дочерние классы:
- WalletHandler - обрабатывает запросы с путем /wallet
  * Методы:
  * ServeHTTP() - проводит базовые проверки на соответсвие спецификации запроса, в зависимости от http метода направляет другим обработчикам:
  * GetMethod() - извлекают аргументы из тела пакета, передают функциям уровня логики, запаковывают в json ответы 
  * PutMethod()
- TransactionHandler - обрабатывает запросы с путем /transaction, по назначению аналогичен Wallet Handler
  * ServeHTTP() 
  * PutMethod()
  * UpdateMethod()
  * DeleteMethod()
  
#### businessLayer
Данные для инициализации: базовая структура пакета sqllayer для вызова методов уровня доступа к БД
Задача: Реализовывать логику приложения
+ Основной класс:
+ Business - Содержит ссылку на базовый класс уровня базы данных
    * Методы:
    * TopUpWallet() - Получает данные о зачислении пользователю денег. Если пользователь не существует, создает нового, если существует, прибавляет к его балансу.
    * GetUserBalance() - Получает идентификатор пользователя, возвращает сведения о балансе пользователя
    * PlaceTransation() - Получает данные для создания нового заказа - создает заказ, списывает сумму у пользователя
    * CommitTransaction() - Получает сигнал признать заказ исполненым, фиксирует дату и время смены статуса
    * CancelTransaction() - Получает сигнал отменить заказ, возвращает зарезервированную сумму.
  


#### sqllayer
Данные для инициализации: параметры подключения к базе данных.
Задача: Работа с базой данных, запаковка полученных из базы данных сведений в структуры.
+ Основной класс: 
- Database
    * Методы: 
    * MetUserById()
    * UpdateUser()
    * CreateUser()
    * GetOrderById()
    * CreateOrder()
    * GetOrderStateById()
    * CreateOrderState()
    * GetServiceById()
    * CreateService()
    
## Трудности
Больше всего вопросов возникло с пакетом "net/http". Поэтому реализация обработчиков отняла много времени. Большая часть статей в интернете по работе с обработчиками эндпоинтов ссылалась на нестандартные маршрутизаторы и фреймворки при работе с веб-серверами.

## Что не успел реализовать:
+ Не успел профести рефакторинг кода
+ Не успел реализовать логирование при помощи Middleware функций
