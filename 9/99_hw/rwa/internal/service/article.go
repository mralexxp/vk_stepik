package service

import (
	"fmt"
	"github.com/gosimple/slug"
	"math/rand"
	"net/url"
	"rwa/internal/dto"
	"rwa/internal/models"
	"time"
)

type ArticlesStore interface {
	Add(*models.Article) uint64
	Delete(uint64) error
	Get(uint64) (*models.Article, error)
	GetSlugID(string) (uint64, error)
	GetByFilter(*models.ArticleFilter) ([]*models.Article, error)
}

func (s *Service) ArticlesByFilter(query *url.Values) (*dto.ArticlesResponse, error) {
	af := models.NewArticleFilter(query)

	if af.AuthorUsername != "" {
		user, err := s.Users.GetByUsername(af.AuthorUsername)
		if err != nil {
			return nil, err
		}

		af.AuthorID = user.ID
	}

	if af.Limit == 0 {
		af.Limit = DefaultArticleLimit // limit по умолчанию
	}

	filtered, err := s.Articles.GetByFilter(af)
	if err != nil {
		return nil, err
	}

	authors := make([]*models.User, len(filtered))

	for i, article := range filtered {
		authors[i], err = s.Users.GetByID(article.Author)
		if err != nil {
			return nil, err
		}
	}

	response := dto.NewArticlesResponse(filtered, authors)

	return response, nil
}

func (s *Service) CreateArticle(aReq *dto.ArticleRequest, token string) (*dto.ArticleResponse, error) {
	id, ok := s.SM.Check(token)
	if !ok {
		return nil, fmt.Errorf("token is invalid")
	}

	// уникальные окончания
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 4)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	aSlug := slug.Make(aReq.Article.Title) + "-" + string(b)

	article := &models.Article{
		Slug:        aSlug,
		Title:       aReq.Article.Title,
		Description: aReq.Article.Description,
		Body:        aReq.Article.Body,
		TagList:     aReq.Article.TagList,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Author:      id,
		Followers:   make([]uint64, 0),
	}

	_ = s.Articles.Add(article)
	author, err := s.Users.GetByID(id)
	if err != nil {
		return nil, err
	}

	responseData := dto.NewArticleResponseData(article, author)

	return &dto.ArticleResponse{Article: responseData}, nil
}
