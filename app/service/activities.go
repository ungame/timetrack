package service

import (
	"context"
	"database/sql"
	"github.com/ungame/timetrack/app/models"
	"github.com/ungame/timetrack/app/observer"
	"github.com/ungame/timetrack/app/repository"
	"github.com/ungame/timetrack/types"
	"log"
	"time"
)

type ActivitiesService interface {
	StartActivity(ctx context.Context, activity *types.Activity) (*types.Activity, error)
	GetActivity(ctx context.Context, id int64) (*types.Activity, error)
	GetActivities(ctx context.Context) ([]*types.Activity, error)
	FinishActivity(ctx context.Context, id int64) (*types.Activity, error)
	UpdateActivity(ctx context.Context, activity *types.Activity) (*types.Activity, error)
	DeleteActivity(ctx context.Context, id int64) error
}

type activitiesService struct {
	categoriesService    CategoriesService
	activitiesRepository repository.ActivitiesRepository
	obs                  observer.Observer
}

func NewActivitiesService(
	categoriesService CategoriesService,
	activitiesRepository repository.ActivitiesRepository,
	obs observer.Observer,
) ActivitiesService {

	return &activitiesService{
		categoriesService:    categoriesService,
		activitiesRepository: activitiesRepository,
		obs:                  obs,
	}
}

func (s *activitiesService) StartActivity(ctx context.Context, activity *types.Activity) (*types.Activity, error) {
	items, err := s.activitiesRepository.GetByStatus(ctx, models.Started)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	go func(activities []*models.Activity) {
		for _, item := range activities {
			_, err = s.FinishActivity(context.Background(), item.ID)
			if err != nil {
				log.Printf("unable to finish current activity: ID=%d, Error=%s\n", item.ID, err.Error())
			}
		}
	}(items)

	newActivity := &models.Activity{
		CategoryID:  activity.CategoryID,
		Description: activity.Description,
		Status:      models.Started,
	}

	newActivity.SetStartedAt(time.Now())
	newActivity.SetUpdatedAt(time.Now())

	newActivity, err = s.activitiesRepository.Create(ctx, newActivity)
	if err != nil {
		return nil, err
	}

	log.Printf("activity started: ID=%d\n", newActivity.ID)

	category, err := s.categoriesService.GetCategory(ctx, newActivity.CategoryID)
	if err == nil {
		s.obs.Count("started", category.Name)
	} else {
		log.Println("unable to get category:", err.Error())
	}

	return newActivity.Type(), nil
}

func (s *activitiesService) GetActivity(ctx context.Context, id int64) (*types.Activity, error) {
	activity, err := s.activitiesRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return activity.Type(), nil
}

func (s *activitiesService) GetActivities(ctx context.Context) ([]*types.Activity, error) {
	items, err := s.activitiesRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	activities := make([]*types.Activity, 0, len(items))
	for _, item := range items {
		activities = append(activities, item.Type())
	}
	return activities, nil
}

func (s *activitiesService) FinishActivity(ctx context.Context, id int64) (*types.Activity, error) {
	activity, err := s.activitiesRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	activity.Status = models.Finished
	activity.SetFinishedAt(time.Now())
	activity.SetUpdatedAt(time.Now())
	_, err = s.activitiesRepository.Update(ctx, activity)
	if err != nil {
		return nil, err
	}

	log.Printf("activity finished: ID=%d\n", id)

	category, err := s.categoriesService.GetCategory(ctx, activity.CategoryID)
	if err == nil {
		s.obs.Count("finished", category.Name)
		s.obs.DurationOf("finished", category.Name, activity.StartedAt.Time)
	} else {
		log.Println("unable to get category:", err.Error())
	}

	return activity.Type(), nil
}

func (s *activitiesService) UpdateActivity(ctx context.Context, activity *types.Activity) (*types.Activity, error) {
	existing, err := s.activitiesRepository.Get(ctx, activity.ID)
	if err != nil {
		return nil, err
	}
	existing.CategoryID = activity.CategoryID
	existing.Description = activity.Description
	existing.SetUpdatedAt(time.Now())
	_, err = s.activitiesRepository.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	log.Printf("activity updated: ID=%d\n", existing.ID)

	category, err := s.categoriesService.GetCategory(ctx, existing.CategoryID)
	if err == nil {
		s.obs.Count("updated", category.Name)
	} else {
		log.Println("unable to get category:", err.Error())
	}

	return existing.Type(), nil
}

func (s *activitiesService) DeleteActivity(ctx context.Context, id int64) error {
	existing, err := s.activitiesRepository.Get(ctx, id)
	if err != nil {
		return err
	}
	rows, err := s.activitiesRepository.Delete(ctx, existing.ID)
	if err != nil {
		return err
	}
	if rows > 0 {
		category, err := s.categoriesService.GetCategory(ctx, existing.CategoryID)
		if err == nil {
			s.obs.Count("deleted", category.Name)
		} else {
			log.Println("unable to get category:", err.Error())
		}
	}
	return nil
}
