package models

import (
	"net/url"
	"time"
)

type Article struct {
	ID          uint64
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	TagList     []string  `json:"tagList"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Author      uint64    `json:"author"`
	Followers   []uint64  `json:"followers"`
}

type ArticleFilter struct {
	Tag            string `schema:"tag"`
	AuthorUsername string `schema:"author"`
	AuthorID       uint64
	Limit          int `schema:"limit"`
	Offset         int `schema:"offset"`
}

func NewArticleFilter(values *url.Values) *ArticleFilter {
	var (
		limitI  interface{}
		offsetI interface{}
	)

	filter := &ArticleFilter{
		Tag:            values.Get("tag"),
		AuthorUsername: values.Get("author"),
	}

	limitI = values.Get("limit")

	if limit, ok := limitI.(int); ok {
		filter.Limit = limit
	}

	offsetI = values.Get("offset")

	if offset, ok := offsetI.(int); ok {
		filter.Offset = offset
	}

	return filter
}
