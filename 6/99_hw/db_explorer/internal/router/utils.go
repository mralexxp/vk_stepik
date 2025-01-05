package router

import (
	"db_explorer/internal/models"
	"net/http"
	"strconv"
)

func GetParams(r *http.Request) *models.QueryParams {
	params := models.NewQueryParams()

	if limit, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		params.Limit = limit
	}

	if offset, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
		params.Offset = offset
	}

	return params
}
