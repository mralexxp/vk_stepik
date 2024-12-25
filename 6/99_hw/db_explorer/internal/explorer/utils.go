package explorer

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func IsExistRow(db *sql.DB, table string, field string, value interface{}) (bool, error) {
	var scanner interface{}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s=? LIMIT 1", field, table, field)

	row := db.QueryRow(query, value)

	err := row.Scan(&scanner)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// Возвращает True и err == nil, если запись найдена
// Возвращает False и err == nil, если запись не найдена
// err != nil при любой ошибке
func (e *Explorer) IsValidField(tableName string, field string, value interface{}) (found bool, err error) {
	const OP = "Explorer.IsValidField"

	if column, ok := e.Struct[tableName][field]; ok {
		switch v := value.(type) {
		case string:
			if strings.Contains(column.Type, "varchar") || strings.Contains(column.Type, "text") {
				found, err := IsExistRow(e.DB, tableName, field, v)
				if err != nil {
					return false, err
				}

				return found, nil
			}

			return false, fmt.Errorf("field %s value %v is not %s type", tableName, value, column.Type)
		// TODO: Обязательно протестировать с точными дробными числами
		// JSON по умолчанию кидает цифры в float, потому необходимо разделить float и int
		// TODO: Если мы попытаемся записать точное дробное число (напр. 4.00), то оно определится как int
		// JSON даже отправленный 4.00 бросает в интерфейс как "4"
		case float64:
			if v == float64(int(v)) {
				if strings.Contains(column.Type, "int") {
					found, err := IsExistRow(e.DB, tableName, field, v)
					if err != nil {
						return false, err
					}

					return found, nil
				}

				return false, fmt.Errorf("field %s value %v is not %s type", tableName, value, column.Type)
			} else if v != float64(int(v)) {
				if strings.Contains(column.Type, "float") {
					found, err := IsExistRow(e.DB, tableName, field, v)
					if err != nil {
						return false, err
					}

					return found, nil
				}
				return false, fmt.Errorf("field %s value %v is not %s type", tableName, value, column.Type)

			} else {
				return false, fmt.Errorf(OP + ": unkown error")
			}
		// TODO: Тест с null
		case nil:
			if column.Null == "YES" {
				found, err := IsExistRow(e.DB, tableName, field, v)
				if err != nil {
					return false, err
				}

				return found, nil
			}

		// На всякий случай оставим int, если вдруг источником данных будет что-то кроме json
		case int:
			if strings.Contains(column.Type, "int") {
				found, err := IsExistRow(e.DB, tableName, field, v)
				if err != nil {
					return false, err
				}

				return found, nil
			}

			return false, fmt.Errorf("field %s value %v is not %s type", tableName, value, column.Type)

		}
	}
	return false, nil
}

func ValidateData(e *Explorer, tableName string, data map[string]interface{}) error {
	// data - словарь с полученными данными из запроса
	// k - имя поля, которое необходимо вписать
	// v - значение этого поля
	// column - содержит структуру Column, если таковое найдено в структуре БД
	/*
		Field = {string} "id"
		Type = {string} "int(11)"
		Null = {string} "NO"
		Key = {string} "PRI"
		Default = {interface{}} nil
		Increment = {string} "auto_increment"
	*/
	for k, v := range data {
		_ = v
		// Проверка на существование
		if column, ok := e.Struct[tableName][k]; ok {
			// Проверка на существование записи
			if ok, err := e.IsValidField(tableName, column.Field, v); err != nil {
				// ERROR
				return err
			} else if ok {
				// Если PRIMARY и запись существует, то обновлять нельзя
				if column.Key == "PRI" {
					return fmt.Errorf("field '%s' value '%v' is primary & exists", k, v)
				}

				fmt.Println("FOUND")
			} else {
				fmt.Println("NOT FOUND")
			}
			// Вся логика проверки

			fmt.Println(column)
		} else {
			return fmt.Errorf("field %s not found from table %s", k, tableName)
		}
	}

	return nil
}
