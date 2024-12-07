package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
)

const (
	fileName = "dataset.xml"
)

func main() {
	http.HandleFunc("/", SearchServer)

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		panic(err)
	}
}

// Handler
func SearchServer(w http.ResponseWriter, r *http.Request) {
	sr, err := parseRequest(r)
	if err != nil {
		responseError(w, err)
		return
	}

	users, err := Search(sr)
	if err != nil {
		responseError(w, NewErr(err, http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = NewErr(json.NewEncoder(w).Encode(*users), http.StatusInternalServerError)
	if err != nil {
		// TODO: заменить панику: не падать во время обрыва
		panic("write response error: " + err.Error())
	}
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
			fmt.Errorf("unknown order_field: %s", r.URL.Query().Get("order_field")),
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
		if limit < 0 {
			return nil, NewErr(
				fmt.Errorf("limit must be positive, recieve: %d", limit),
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
