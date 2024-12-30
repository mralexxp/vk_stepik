package explorer

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (e *Explorer) IsExistRowFromPrimary(table string, id int) (bool, error) {
	var priKey string

	for k := range e.Struct[table] {
		if e.Struct[table][k].Key == "PRI" {
			priKey = k
			break
		}
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, priKey)
	row := e.DB.QueryRow(query, id)

	var scanner interface{}

	err := row.Scan(&scanner)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (e *Explorer) IsValidField(tableName string, field string, value interface{}) (bool, error) {
	const OP = "Explorer.IsValidField"

	// TODO: Убрать проверку наличия поля, так как проверяем наличие в родительской функции
	if column, ok := e.Struct[tableName][field]; ok {
		switch v := value.(type) {
		case string:
			if strings.Contains(column.Type, "varchar") || strings.Contains(column.Type, "text") {
				return true, nil
			}

			return false, fmt.Errorf("field %s have invalid type", field)
		// TODO: Обязательно протестировать с точными дробными числами
		// JSON по умолчанию кидает цифры в float, потому необходимо разделить float и int
		// TODO: Если мы попытаемся записать точное дробное число (напр. 4.00), то оно определится как int
		// JSON даже отправленный 4.00 бросает в интерфейс как "4"
		case float64:
			if v == float64(int(v)) {
				if strings.Contains(column.Type, "int") {
					return true, nil
				}

				return false, fmt.Errorf("field %s have invalid type", field)
			} else if v != float64(int(v)) {
				if strings.Contains(column.Type, "float") {
					return true, nil
				}

				return false, fmt.Errorf("field %s have invalid type", field)
			} else {
				return false, fmt.Errorf(OP + ": unkown error")
			}
		// TODO: Тест с null
		case nil:
			if column.Null == "YES" {
				return true, nil
			}

			return false, fmt.Errorf("field %s have invalid type", field)

		// На всякий случай оставим int, если вдруг источником данных будет что-то кроме json
		case int:
			if strings.Contains(column.Type, "int") {
				return true, nil
			}

			return false, fmt.Errorf("field %s have invalid type", field)
		default:
			return false, fmt.Errorf("%s: unkown error", OP)
		}
	} else {
		return false, fmt.Errorf("field %s have invalid type", field)
	}
}

func emptyValue(data *map[string]interface{}, column map[string]Column) {
	for k := range column {
		if _, ok := (*data)[k]; !ok && column[k].Null != "YES" {
			switch {
			case strings.Contains(column[k].Type, "text") || strings.Contains(column[k].Type, "varchar"):
				(*data)[k] = ""
			case strings.Contains(column[k].Type, "int"):
				(*data)[k] = 0
			default:
				(*data)[k] = ""
			}
		}
	}
}

func (e *Explorer) InsertConstructor(tableName string, data map[string]interface{}) (query string, placeholders []interface{}, err error) {
	placeholders = make([]interface{}, 0)
	values := make([]string, 0)
	columns := make([]string, 0)

	for k, v := range data {
		values = append(values, "?")
		columns = append(columns, k)
		placeholders = append(placeholders, v)
	}
	query = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		tableName,
		strings.Join(columns, ","),
		strings.Join(values, ","),
	)

	return query, placeholders, nil
}

func (e *Explorer) UpdateConstructor(tableName string, id int, data map[string]interface{}) (query string, placeholders []interface{}, err error) {
	var priKey string

	for k := range e.Struct[tableName] {
		if e.Struct[tableName][k].Key == "PRI" {
			priKey = k
			break
		}
	}

	query = "UPDATE " + tableName + " SET "
	i := 0
	for k, v := range data {
		query += fmt.Sprintf("%s = ?, ", k)
		if v == "" {
			placeholders = append(placeholders, nil)
		} else {
			placeholders = append(placeholders, v)
		}

		if len(data)-1 == i {
			query = query[:len(query)-2]
		}

		i++
	}

	query += " WHERE " + priKey + " = ?"
	placeholders = append(placeholders, id)

	return query, placeholders, nil

}
