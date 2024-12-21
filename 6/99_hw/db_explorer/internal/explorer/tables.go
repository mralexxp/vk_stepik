package explorer

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
