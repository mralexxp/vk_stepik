package articles

import (
	"rwa/internal/models"
	"sort"
)

type SortedArticles struct {
	articles []*models.Article
	dateIdx  map[uint64]int
}

func NewSortedArticles() *SortedArticles {
	return &SortedArticles{
		articles: make([]*models.Article, 0),
		dateIdx:  make(map[uint64]int),
	}
}

func (sa *SortedArticles) Add(article *models.Article) {
	pos := sort.Search(len(sa.articles), func(i int) bool {
		return sa.articles[i].CreatedAt.After(article.CreatedAt)
	})

	sa.articles = append(sa.articles, nil)
	copy(sa.articles[pos+1:], sa.articles[pos:])
	sa.articles[pos] = article

	for i := pos; i < len(sa.articles); i++ {
		sa.dateIdx[sa.articles[i].ID] = i
	}
}

func (sa *SortedArticles) Delete(articleID uint64) {
	if pos, exists := sa.dateIdx[articleID]; exists {
		sa.articles = append(sa.articles[:pos], sa.articles[pos+1:]...)
		delete(sa.dateIdx, articleID)

		for i := pos; i < len(sa.articles); i++ {
			sa.dateIdx[sa.articles[i].ID] = i
		}
	}
}

func (sa *SortedArticles) GetArticles(offset, limit int) []*models.Article {
	if offset >= len(sa.articles) {
		return []*models.Article{}
	}

	end := offset + limit
	if end > len(sa.articles) {
		end = len(sa.articles)
	}

	return sa.articles[offset:end]
}

// Дополнительный метод для фильтрации статей по ID
func (sa *SortedArticles) FilterByIDs(ids map[uint64]struct{}) []*models.Article {
	if len(ids) == 0 {
		return sa.articles
	}

	filtered := make([]*models.Article, 0, len(ids))
	for _, article := range sa.articles {
		if _, exists := ids[article.ID]; exists {
			filtered = append(filtered, article)
		}
	}
	return filtered
}
