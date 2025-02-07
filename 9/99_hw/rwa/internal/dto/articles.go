package dto

import (
	"rwa/internal/models"
	"time"
)

type ArticleRequest struct {
	Article *ArticleRequestData `json:"article"`
}

type ArticleResponse struct {
	Article *ArticleResponseData `json:"article"`
}

type ArticlesResponse struct {
	Articles      []*ArticleResponseData `json:"articles"`
	ArticlesCount int                    `json:"articlesCount"`
}

type ArticleRequestData struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

type ArticleResponseData struct {
	ID          uint64                 `json:"id"`
	Slug        string                 `json:"slug,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Body        string                 `json:"body,omitempty"`
	TagList     []string               `json:"tagList,omitempty"`
	CreatedAt   string                 `json:"createdAt,omitempty"`
	UpdatedAt   string                 `json:"updatedAt,omitempty"`
	Author      *ArticleAuthorResponse `json:"author,omitempty"`
}

type ArticleAuthorResponse struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func NewArticlesResponse(articles []*models.Article, authors []*models.User) *ArticlesResponse {
	responseData := make([]*ArticleResponseData, len(articles))

	for i, article := range articles {
		responseData[i] = &ArticleResponseData{
			ID:          article.ID,
			Slug:        article.Slug,
			Title:       article.Title,
			Description: article.Description,
			Body:        article.Body,
			TagList:     article.TagList,
			CreatedAt:   article.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   article.UpdatedAt.Format(time.RFC3339),
			Author: &ArticleAuthorResponse{
				Username: authors[i].Username,
				Bio:      authors[i].Bio,
				Image:    authors[i].Image,
			},
		}
	}

	return &ArticlesResponse{
		Articles:      responseData,
		ArticlesCount: len(articles),
	}
}

func NewArticleAuthorResponse(author *models.User) *ArticleAuthorResponse {
	return &ArticleAuthorResponse{
		Username: author.Username,
		Bio:      author.Bio,
		Image:    author.Image,
	}
}

func NewArticleResponseData(article *models.Article, author *models.User) *ArticleResponseData {
	return &ArticleResponseData{
		Slug:        article.Slug,
		Title:       article.Title,
		Description: article.Description,
		Body:        article.Body,
		TagList:     article.TagList,
		CreatedAt:   article.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   article.UpdatedAt.Format(time.RFC3339),
		Author:      NewArticleAuthorResponse(author),
	}
}
