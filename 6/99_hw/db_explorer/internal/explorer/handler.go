package explorer

import "encoding/json"

// Возвращает готовй к отправке слайс байт
// Если есть ошибки вернет пустой слайс и ошибку
func (e *Explorer) ShowTables() ([]byte, error) {
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

	response, err := json.Marshal(tables)
	if err != nil {
		return nil, err
	}

	return response, nil
}
