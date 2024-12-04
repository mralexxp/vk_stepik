package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Usr struct {
	Id        int `xml:"id"`
	Name      string
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type Request struct {
	Query string `xml:"query"` // Ищем по полям записи `Name` и `About` просто подстроку, без регулярок.
	// `Name` - это first_name + last_name из xml.
	// Поле пустое-возвращаем все записи (поиск пустой подстроки всегда возвращает true) делаем только логику сортировки
	OrderField string `xml:"order_field"` // По какому полю сортировать.
	// Работает по полям `Id`, `Age`, `Name`, если пустой - то сортируем по `Name`
	//если что-то другое - SearchServer ругается ошибкой.
	OrderBy string `xml:"order_by"` // Направление сортировки (как есть, по убыванию, по возрастанию),
	// в client.go есть соответствующие константы
	Limit  int `xml:"limit"`  // Сколько записей вернуть
	Offset int `xml:"offset"` // Начиная с какой записи вернуть (сколько пропустить с начала)
}

const (
	fileName = "dataset.xml"
)

func main() {

	realUsers := SearchServer("", "Name", OrderByAsc, 0, 0)
	for _, realUser := range realUsers {
		fmt.Println(realUser)
	}

}

// В отдельном хендлере
func SearchServer(query string, orderField string, orderBy int, limit int, offset int) []Usr {
	// Если поле для сортировки задано неверно
	if orderField != "Name" && orderField != "Id" && orderField != "Age" && orderField != "" {
		return []Usr{}
	}

	type fileStruct struct {
		Root []Usr `xml:"row"`
	}

	users := make([]Usr, 0)
	found := 0

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	dataXml, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	fs := fileStruct{}
	err = xml.Unmarshal(dataXml, &fs)
	if err != nil {
		panic(err)
	}

LOOP:
	for _, user := range fs.Root {
		user.Name = user.FirstName + " " + user.LastName
		if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
			found++

			if found == limit+offset+1 && limit != 0 {
				break LOOP
			}

			if offset >= found {
				continue
			}

			user.About = strings.Replace(user.About, "\n", "", 1)
			users = append(users, user)
		}
	}

	switch orderBy {
	case OrderByAsc:
		err = SortSlices(users, orderField)
		if err != nil {
			panic(err)
		}

		Reverse(users)

		return users
	case OrderByDesc:
		err = SortSlices(users, orderField)
		if err != nil {
			panic(err)
		}

		return users
	default:
		return users
	}
}

func SortSlices(users []Usr, orderField string) error {
	switch orderField {
	case "Id":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Id < users[j].Id
		})
	case "Age":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Age < users[j].Age
		})
	case "Name":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name < users[j].Name
		})
	case "":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name < users[j].Name
		})
	default:
		return fmt.Errorf("invalid orderField: %s", orderField)
	}

	return nil
}

func Reverse(u []Usr) {
	for i, j := 0, len(u)-1; i < len(u)/2; i, j = i+1, j-1 {
		u[i], u[j] = u[j], u[i]
	}
}
