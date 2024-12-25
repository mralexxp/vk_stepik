package explorer

import (
	"fmt"
	"strings"
)

// Возвращает готовй к отправке слайс байт
// Если есть ошибки вернет пустой слайс и ошибку
func (e *Explorer) ShowTables() ([]string, error) {
	query := "SHOW TABLES;"
	rows, err := e.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		table := ""
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	// TODO: маршаллингом занимаемся в ручке или пишем io.Writer?
	return tables, nil
}

func (e *Explorer) ShowTable(tableName string, params map[string]int) ([]map[string]interface{}, error) {
	// защита от инъекций, так как database/sql не поддерживает placeholders для имен таблиц
	// Не проверяем на инъекции, так как пользователь не попадет в эту функцию при попытке инъекции
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
			switch v := val.(type) {
			case []byte:
				tempResult[columns[i]] = string(v)
			case nil:
				tempResult[columns[i]] = nil
			}

		}
		result = append(result, tempResult)
	}

	// ВЕРНУТЬ []map[<field_name>]<value>
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

	query := `SELECT * FROM` + ` ` + tableName + ` WHERE id=? LIMIT 1`

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
			}

		}
	}

	return result, nil

}

func (e *Explorer) PutTuple(tableName string, data map[string]interface{}) (int, error) {
	// TODO: VALIDATE Increment:
	for k, v := range data {
		if stct, ok := e.Struct[tableName][k]; ok {
			if stct.Increment != "" && v != "" {
				delete(data, k)
			}
		} else {
			return -1, fmt.Errorf("field %s not found", tableName)
		}
	}

	query, placeholders, err := InsertConstructor(tableName, data)
	if err != nil {
		return 0, err
	}

	exec, err := e.DB.Exec(query, placeholders...)
	if err != nil {
		return 0, err
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (e *Explorer) UpdateTuple(tableName string, id int, data map[string]interface{}) (int, error) {
	//// TODO: Validate
	err := ValidateData(e, tableName, data)
	if err != nil {
		return 0, err
	}

	query, placeholders, err := UpdateConstructor(tableName, id, data)
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
	query := fmt.Sprintf("DELETE FROM `%s` WHERE id=?", tableName)
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

func InsertConstructor(tableName string, data map[string]interface{}) (query string, placeholders []interface{}, err error) {
	// INSERT INTO <tablename>(<columns>) VALUES (<placeholders>)
	placeholders = make([]interface{}, 0)
	values := make([]string, 0)
	columns := make([]string, 0)

	for k, v := range data {
		values = append(values, "?")
		columns = append(columns, k)
		//if v == "" {
		//	placeholders = append(placeholders, nil)
		//	continue
		//}
		placeholders = append(placeholders, v)
	}
	query = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		tableName,
		strings.Join(columns, ","),
		strings.Join(values, ","),
	)

	return query, placeholders, nil
}

func UpdateConstructor(tableName string, id int, data map[string]interface{}) (query string, placeholders []interface{}, err error) {
	// UPDATE <tablename>
	// SET <field> = <value>,
	// 	   <field> = <value>
	// WHERE id = <id>

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

	query += " WHERE id = ?"
	placeholders = append(placeholders, id)

	return query, placeholders, nil
}
