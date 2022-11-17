package service

import (
	"context"
	"github.com/ungame/timetrack/app/models"
	"github.com/ungame/timetrack/app/repository"
	"github.com/ungame/timetrack/types"
	"log"
	"sync"
)

type CategoriesService interface {
	GetCategory(ctx context.Context, id int64) (*types.Category, error)
	GetCategories(ctx context.Context) ([]*types.Category, error)
}

type categoriesService struct {
	categoriesRepository repository.CategoriesRepository
	mutex                *sync.RWMutex
	cache                map[int64]*models.Category
}

func NewCategoriesService(categoriesRepository repository.CategoriesRepository) CategoriesService {
	svc := &categoriesService{
		categoriesRepository: categoriesRepository,
	}
	svc.load()
	return svc
}

func (s *categoriesService) load() {
	categories, err := s.categoriesRepository.GetAll(context.Background())
	if err != nil {
		log.Panicln("unable to cache categories:", err.Error())
	}
	s.mutex = &sync.RWMutex{}
	s.cache = make(map[int64]*models.Category, len(categories))
	for _, category := range categories {
		s.cache[category.ID] = category
	}
}

func (s *categoriesService) GetCategory(ctx context.Context, id int64) (*types.Category, error) {
	s.mutex.RLock()
	if category, ok := s.cache[id]; ok {
		s.mutex.RUnlock()
		return category.Type(), nil
	}

	category, err := s.categoriesRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	s.mutex.Lock()
	s.cache[category.ID] = category
	s.mutex.Unlock()

	return category.Type(), nil
}

func (s *categoriesService) GetCategories(ctx context.Context) ([]*types.Category, error) {
	items, err := s.categoriesRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	categories := make([]*types.Category, 0, len(items))
	for _, item := range items {
		categories = append(categories, item.Type())
	}
	return categories, nil
}
