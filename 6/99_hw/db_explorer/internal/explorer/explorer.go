package explorer

import "database/sql"

type Explorer struct {
	DB     *sql.DB
	Struct map[string]map[string]Column
}

type Column struct {
	Field     string      `sql:"Field"`
	Type      string      `sql:"Type"`
	Null      string      `sql:"Null"`
	Key       string      `sql:"Key"`
	Default   interface{} `sql:"Default"` // Пока не придумал как использовать
	Increment string      `sql:"increment"`
}

func (e *Explorer) InitDBStruct() {
	// Получаем список таблиц
	rows, err := e.DB.Query("SHOW TABLES;")
	if err != nil {
		// Без возврата ошибки, так как паника вызывается при запуске программы
		// потому не доставит проблем во время работы
		panic(err)
	}

	defer rows.Close()

	tables := make([]string, 0)
	table := ""

	for rows.Next() {
		err := rows.Scan(&table)
		if err != nil {
			panic(err)
		}

		tables = append(tables, table)
	}

	rows.Close()

	e.Struct = make(map[string]map[string]Column)

	// Получаем поля таблиц
	for _, table := range tables {
		// !!! Обязательная валидация имени таблицы (var table) при вызове от пользователя !!!
		e.Struct[table] = make(map[string]Column)
		rows, err = e.DB.Query("SHOW COLUMNS FROM " + table)
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			field := Column{}

			err := rows.Scan(
				&field.Field,
				&field.Type,
				&field.Null,
				&field.Key,
				&field.Default,
				&field.Increment,
			)

			if err != nil {
				panic(err)
			}

			e.Struct[table][field.Field] = field
		}

		rows.Close()
	}

}