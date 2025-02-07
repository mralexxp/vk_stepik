package service

import (
	"fmt"
	"github.com/gosimple/slug"
	"net/url"
	"rwa/internal/dto"
	"rwa/internal/models"
	"time"
)

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

	article := &models.Article{
		Slug:        slug.Make(aReq.Article.Title),
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
