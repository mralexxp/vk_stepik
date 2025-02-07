package articles

import "rwa/internal/models"

type Index struct {
	slug    map[string]uint64
	tags    map[string]map[uint64]struct{}
	authors map[uint64]map[uint64]struct{} // authors[userID][articleID]
	sorted  *SortedArticles
}

func NewIndex() *Index {
	return &Index{
		slug:    make(map[string]uint64),
		tags:    make(map[string]map[uint64]struct{}),
		authors: make(map[uint64]map[uint64]struct{}),
		sorted:  NewSortedArticles(),
	}
}

func (i *Index) Add(a *models.Article) {
	i.slug[a.Slug] = a.ID

	for _, tag := range a.TagList {
		if _, ok := i.tags[tag]; !ok {
			i.tags[tag] = make(map[uint64]struct{})
		}

		i.tags[tag][a.ID] = struct{}{}
	}

	if _, ok := i.authors[a.Author]; !ok {
		i.authors[a.Author] = make(map[uint64]struct{})
	}

	i.authors[a.Author][a.ID] = struct{}{}

	// sorted
	i.sorted.Add(a)
}

func (i *Index) Delete(a *models.Article) {
	delete(i.slug, a.Slug)

	for _, tag := range a.TagList {
		if _, ok := i.tags[tag]; ok {
			delete(i.tags[tag], a.ID)
		}
	}

	delete(i.authors[a.Author], a.ID)

	i.sorted.Delete(a.ID)
}

func (i *Index) GetBySlug(slug string) uint64 {
	if id, ok := i.slug[slug]; ok {
		return id
	}

	return 0
}

func (i *Index) GetByAuthor(author uint64) []uint64 {
	if a, ok := i.authors[author]; ok {
		result := make([]uint64, len(a))
		for k := range a {
			result = append(result, k)
		}

		return result
	}

	return nil
}

func (i *Index) GetByTag(tag string) []uint64 {
	if a, ok := i.tags[tag]; ok {
		result := make([]uint64, len(a))
		for k := range a {
			result = append(result, k)
		}

		return result
	}

	return nil
}
