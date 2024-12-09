package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
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

// Основная логика
func Search(s *SearchRequest) (*[]User, CustomError) {
	users := make([]User, 0)

	allUsers, err := LoadUsers(useXML)
	if err != nil {
		return nil, NewErr(err, http.StatusInternalServerError)
	}

	for _, user := range *allUsers {
		user.Name = user.FirstName + " " + user.LastName
		if strings.Contains(user.Name, s.Query) || strings.Contains(user.About, s.Query) {
			u := User{
				Id:     user.Id,
				Name:   user.Name,
				Age:    user.Age,
				About:  strings.Replace(user.About, "\n", "", 1),
				Gender: user.Gender,
			}

			users = append(users, u)
		}
	}

	switch s.OrderBy {
	case OrderByAsc:
		err = SortSlices(users, s.OrderField)
		if err != nil {
			return nil, NewErr(err, http.StatusInternalServerError)
		}

		Reverse(users)
	case OrderByDesc:
		err = SortSlices(users, s.OrderField)
		if err != nil {
			return nil, NewErr(err, http.StatusInternalServerError)
		}
	case OrderByAsIs:
	default:
		return nil, NewErr(err, http.StatusBadRequest)
	}

	if s.Offset > len(users) {
		return &[]User{}, nil
	}

	users = users[s.Offset:]

	if s.Limit >= len(users) || s.Limit == 0 {
		return &users, NewErr(err, http.StatusOK)
	}

	users = users[:s.Limit]
	return &users, NewErr(err, http.StatusOK)
}

func SortSlices(users []User, orderField string) CustomError {
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
		return NewErr(fmt.Errorf("unknown orderField: %s", orderField), http.StatusInternalServerError)
	}

	return nil
}

func Reverse(u []User) {
	for i, j := 0, len(u)-1; i < len(u)/2; i, j = i+1, j-1 {
		u[i], u[j] = u[j], u[i]
	}
}

func parseRequest(r *http.Request) (*SearchRequest, CustomError) {
	sr := SearchRequest{}

	sr.Query = r.URL.Query().Get("query")

	orderField := r.URL.Query().Get("order_field")
	if orderField != "Name" && orderField != "Id" && orderField != "Age" && orderField != "" {
		return nil, NewErr(
			//fmt.Errorf("unknown order_field: %s", r.URL.Query().Get("order_field")),
			fmt.Errorf("ErrorBadOrderField"),
			http.StatusBadRequest,
		)
	}

	sr.OrderField = orderField

	if orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by")); err != nil {
		return nil, NewErr(
			fmt.Errorf("unknown order_by: %v", r.URL.Query().Get("order_by")),
			http.StatusBadRequest,
		)
	} else {
		if orderBy > 1 || orderBy < -1 {
			return nil, NewErr(
				fmt.Errorf("unknown order_by: %v", r.URL.Query().Get("order_by")),
				http.StatusBadRequest,
			)
		}
		sr.OrderBy = orderBy
	}

	if limit, err := strconv.Atoi(r.URL.Query().Get("limit")); err != nil {
		return nil, NewErr(
			fmt.Errorf("unknown limit: %v", r.URL.Query().Get("limit")),
			http.StatusBadRequest,
		)
	} else {
		if limit <= 0 {
			return nil, NewErr(
				fmt.Errorf("expected limit > 0, recieve: %d", limit),
				http.StatusBadRequest,
			)
		}
		sr.Limit = limit
	}

	if offset, err := strconv.Atoi(r.URL.Query().Get("offset")); err != nil {
		return nil, NewErr(
			fmt.Errorf("unknown offset: %v", r.URL.Query().Get("offset")),
			http.StatusBadRequest,
		)
	} else {
		if offset < 0 {
			return nil, NewErr(
				fmt.Errorf("offset must be positive, recieve: %d", offset),
				http.StatusBadRequest,
			)
		}
		sr.Offset = offset
	}

	return &sr, nil
}

func LoadUsers(xmlFilename string) (*[]Usr, CustomError) {
	file, err := os.Open(xmlFilename)
	if err != nil {
		return nil, NewErr(err, http.StatusInternalServerError)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	dataXml, err := io.ReadAll(file)
	if err != nil {
		return nil, NewErr(err, http.StatusInternalServerError)
	}

	fs := struct {
		Root []Usr `xml:"row"`
	}{}

	err = xml.Unmarshal(dataXml, &fs)
	if err != nil {
		return nil, NewErr(err, http.StatusInternalServerError)
	}

	return &fs.Root, nil
}
