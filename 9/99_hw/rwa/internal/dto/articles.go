package dto

import "time"

type ArticleRequest struct {
	Article *ArticleRequestData `json:"article"`
}

type ArticleResponse struct {
	Article *ArticleResponseData `json:"article"`
}

type ArticlesResponse struct {
	Articles []*ArticleResponseData `json:"articles"`
}

type ArticleResponseData struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	TagList     []string  `json:"tagList"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Author      uint64    `json:"author"`
}

type ArticleRequestData struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}
