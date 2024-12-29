package router

import (
	"net/http"
	"strconv"
)

func GetParams(r *http.Request) map[string]int {
	params := make(map[string]int)

	var limit int
	var offset int
	var err error

	if limit, err = strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		params["limit"] = limit
	} else {
		params["limit"] = 5
	}

	if offset, err = strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
		params["offset"] = offset
	} else {
		params["offset"] = 0
	}

	return params
}
