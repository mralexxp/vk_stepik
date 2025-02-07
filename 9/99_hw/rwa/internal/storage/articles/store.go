package articles

import (
	"fmt"
	"rwa/internal/models"
	"sync"
)

type Store struct {
	DB    map[uint64]*models.Article
	index *Index
	next  uint64
	mu    *sync.Mutex
}

func NewStore() *Store {
	return &Store{
		DB:    make(map[uint64]*models.Article),
		index: NewIndex(),
		next:  1,
		mu:    &sync.Mutex{},
	}
}

func (s *Store) Add(article *models.Article) uint64 {
	article.ID = s.next

	s.DB[article.ID] = article

	s.index.Add(article)

	s.next++

	return article.ID
}

func (s *Store) Delete(id uint64) error {
	if a, exists := s.DB[id]; exists {
		s.index.Delete(a)

		delete(s.DB, id)

		return nil
	}

	return fmt.Errorf("article with ID %d does not exist", id)
}

func (s *Store) Get(id uint64) (*models.Article, error) {
	article, exists := s.DB[id]
	if !exists {
		return nil, fmt.Errorf("article with ID %d does not exist", id)
	}

	return article, nil
}

func (s *Store) GetSlugID(slug string) (uint64, error) {
	id := s.index.GetBySlug(slug)
	if id == 0 {
		return 0, fmt.Errorf("slug %s does not exist", slug)
	}

	return id, nil
}

func (s *Store) GetByFilter(f *models.ArticleFilter) ([]*models.Article, error) {
	if f == nil {
		return []*models.Article{}, fmt.Errorf("filter is nil")
	}

	resultID := make(map[uint64]struct{})
	filters := false

	if f.Tag != "" {
		filters = true
		ids := s.index.GetByTag(f.Tag)
		for _, id := range ids {
			resultID[id] = struct{}{}
		}
	}

	if f.AuthorID != 0 {
		filters = true

		ids := s.index.GetByAuthor(f.AuthorID)
		for _, id := range ids {
			resultID[id] = struct{}{}
		}
	}

	var articles []*models.Article
	if filters {
		articles = s.index.sorted.FilterByIDs(resultID)
	} else {
		articles = s.index.sorted.GetArticles(f.Offset, f.Limit)
	}

	if filters {
		if f.Offset >= len(articles) {
			return []*models.Article{}, nil
		}

		end := f.Offset + f.Limit
		if end > len(articles) {
			end = len(articles)
		}

		articles = articles[f.Offset:end]
	}

	return articles, nil
}
