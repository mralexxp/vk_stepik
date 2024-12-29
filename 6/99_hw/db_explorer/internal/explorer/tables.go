package explorer

import (
	"fmt"
	"strconv"
	"strings"
)

func (e *Explorer) GetTables() ([]string, error) {
	query := "SHOW TABLES;"
	rows, err := e.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tables := make([]string, 0)
	var table string

	for rows.Next() {
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (e *Explorer) ShowTable(tableName string, params map[string]int) ([]map[string]interface{}, error) {
	q := fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", tableName, params["limit"], params["offset"])
	rows, err := e.DB.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	values := make([]interface{}, len(columns), len(columns))
	valuesPtr := make([]interface{}, len(values), len(values))

	for i := range values {
		valuesPtr[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(valuesPtr...)
		if err != nil {
			return nil, err
		}

		tempResult := make(map[string]interface{}, 0)
		for i, val := range values {
			if column, ok := e.Struct[tableName][columns[i]]; ok {
				switch {
				case strings.Contains(column.Type, "text") || strings.Contains(column.Type, "varchar"):
					val, ok := val.([]byte)
					if ok {
						tempResult[columns[i]] = string(val)
						break
					}

					if column.Null == "YES" {
						tempResult[columns[i]] = val
					}

				case strings.Contains(column.Type, "int"):
					val, err = strconv.Atoi(string(val.([]byte)))
					if err != nil {
						return nil, err
					}
					tempResult[columns[i]] = val
				case strings.Contains(column.Type, "float"):
					tempResult[columns[i]] = val.(float64)
				default:
					tempResult[columns[i]] = val
				}
			} else {
				// Если нет такой колонки
				return nil, fmt.Errorf("column %s not found", columns[i])
			}
		}
		result = append(result, tempResult)
	}

	return result, nil
}

func (e *Explorer) GetTableStruct(tableName string) ([]string, error) {
	result := make([]string, 0)

	if _, ok := e.Struct[tableName]; ok {
		for key := range e.Struct[tableName] {
			result = append(result, key)
		}
		return result, nil
	}

	return []string{}, fmt.Errorf("table %s not found", tableName)
}

func (e *Explorer) GetTuple(tableName string, id int) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)

	var priKey string

	for k := range e.Struct[tableName] {
		if e.Struct[tableName][k].Key == "PRI" {
			priKey = k
			break
		}
	}

	query := `SELECT * FROM` + ` ` + tableName + ` WHERE ` + priKey + `=? LIMIT 1`

	rows, err := e.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()

	values := make([]interface{}, len(columns), len(columns))
	valuesPtr := make([]interface{}, len(values), len(values))

	for i := range values {
		valuesPtr[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(valuesPtr...)
		if err != nil {
			return nil, err
		}

		for i, val := range values {
			switch v := val.(type) {
			case []byte:
				result[columns[i]] = string(v)
			case nil:
				result[columns[i]] = nil
			default:
				result[columns[i]] = v
			}

		}
	}

	return result, nil

}

func (e *Explorer) PutTuple(tableName string, data map[string]interface{}) (map[string]interface{}, error) {
	priKey := ""

	for k, v := range data {
		if stct, ok := e.Struct[tableName][k]; ok {
			if stct.Increment != "" && v != "" {
				delete(data, k)
			}

			if stct.Key == "PRI" {
				priKey = k
			}
		} else {
			delete(data, k)
			// Почему-то мы должны игнорировать левые поля, хотя вернее вернуть ошибку
			// return -1, fmt.Errorf("field %s not found", k)
		}
	}

	// Не переданные обязательные поля должны заполнять значениями по умолчанию, если поел не может быть null
	emptyValue(&data, e.Struct[tableName])

	query, placeholders, err := e.InsertConstructor(tableName, data)
	if err != nil {
		return nil, err
	}

	stmt, err := e.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	exec, err := stmt.Exec(placeholders...)
	if err != nil {
		return nil, err
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return nil, err
	}

	response := make(map[string]interface{}, 1)
	response[priKey] = int(id)

	return response, nil
}

func (e *Explorer) UpdateTuple(tableName string, id int, data map[string]interface{}) (int, error) {
	// валидируем поля и значения
	for k, v := range data {
		// Проверка существования необходимых полей в нашей котаблице
		column, ok := e.Struct[tableName][k]
		if !ok {
			return -1, fmt.Errorf("field %s not found", k)
		}

		// Если тип полей и данных не совпадает
		validField, err := e.IsValidField(tableName, k, v)
		if err != nil && !validField {
			return -1, err
		}

		// ЕСЛИ ПОЛЕ PRIMARY, то изменять его нельзя!
		if column.Key == "PRI" {
			return -1, fmt.Errorf("field %s have invalid type", k)
		}
	}

	// Проверяем существование записей для изменений
	found, err := e.IsExistRowFromPrimary(tableName, id)
	if err != nil {
		return -1, fmt.Errorf("record %d not found", id)
	}

	if !found {
		return -1, fmt.Errorf("record not found")
	}

	query, placeholders, err := e.UpdateConstructor(tableName, id, data)
	if err != nil {
		return 0, err
	}

	result, err := e.DB.Exec(query, placeholders...)
	if err != nil {
		return 0, err
	}

	updated, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(updated), nil
}

func (e *Explorer) DeleteTuple(tableName string, id int) (int, error) {
	var priKey string

	for k := range e.Struct[tableName] {
		if e.Struct[tableName][k].Key == "PRI" {
			priKey = k
			break
		}
	}

	query := fmt.Sprintf("DELETE FROM `%s` WHERE %s=?", tableName, priKey)
	result, err := e.DB.Exec(query, id)
	if err != nil {
		return 0, err
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(deleted), nil
}

func (e *Explorer) InsertConstructor(tableName string, data map[string]interface{}) (query string, placeholders []interface{}, err error) {
	// INSERT INTO <tablename>(<columns>) VALUES (<placeholders>)
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
		fmt.Println(len(data))
		if len(data)-1 == i {
			query = query[:len(query)-2]
		}

		i++
	}

	query += " WHERE " + priKey + " = ?"
	placeholders = append(placeholders, id)

	return query, placeholders, nil
}
