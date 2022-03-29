# GO-simple-sql-builder
от "Рожденных в СССР" :)

### Простой и элегантный построитель SQL-запросов (MySql 5.1) для языка Go без применения reflect.

Предварительная версия.
1. Одиночные и UNION запросы без оконных функций и прочих прелестей современных версий.
2. Нет метода для формирования условий типа t1.field1 = t2.field2, но есть передача условий в режиме AS-IS (покрывает все недостатки).

1. Позволяет создавать полноценные одиночные SQL запросы, произвольной вложенности рекурсивно.
2. Позволяет собирать из отдельных запросов UNION в любом виде.
3. Позволяет формировать список аргументов к запросу как ...interface{} для последующей передачи в database/sql пакеты, в т.ч. sqlx
4. Есть задел на встроенную реализацию выборок через пакет database/sql без применения дополнительной рефлексии на базе callback-функций.

В планах:
1. Дополнение билдера функциями MySql для упрощения проброса полей.
2. Улучшение в части построения фильтров данных произвольной сложности без привлечения map[string]interface{}
3. Отказ от пакетов sqlx, sql и переход на прямую работу с mySql драйвером, снижение требуемой вложенности рефлексий при работе с sql-запросами.
4. Доработка встроенного sql-mock "до ума"..
