package main

import (
	"net/http"
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

	allUsers, err := LoadUsers(fileName)
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
