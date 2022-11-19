package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ungame/timetrack/app/models"
	"github.com/ungame/timetrack/app/repository"
	"github.com/ungame/timetrack/db"
	"github.com/ungame/timetrack/ioext"
	"github.com/ungame/timetrack/timeext"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const defaultTimeout = time.Second * 3

func Load(startTime time.Time) {
	log.Println("Start loading from: ", startTime.Format(timeext.DateTimeFormat))
	defer log.Println("Stopped.")

	ctx := context.Background()
	conn := db.Lite(db.DefaultFileStorage(), db.NewMigration(db.GetSqliteSeed()))
	defer ioext.Close(conn)

	catRepo := repository.NewCategoriesRepository(conn)

	categories, err := catRepo.GetAll(ctx)
	if err != nil {
		log.Panicln(err)
	}

	actRepo := repository.NewActivitiesRepository(conn)
	defer actRepo.Close()

	err = actRepo.Truncate(ctx)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Truncate table activities successfully.")

	yesterday := time.Now().AddDate(0, 0, -1)
	end := timeext.GetEndOfDayFrom(yesterday, yesterday.Location())

	finishCh := make(chan *models.Activity, 10)
	semaphore := make(chan struct{}, 10)
	startedCounter := 0

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(finish <-chan *models.Activity) {
		defer wg.Done()

		log.Println("Starting loop to finish activities.")

		finishedCounter := 0

		for {
			select {
			case activity := <-finish:
				log.Println("[*] Received Activity:", activity.ID)
				activity.Status = models.Finished
				finishedAt := activity.StartedAt.Time.Add(time.Hour)

				activity.SetUpdatedAt(finishedAt)
				activity.SetFinishedAt(finishedAt)

				timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)

				_, err = actRepo.Update(timeoutCtx, activity)
				if err != nil {
					log.Panicln(err)
				}

				log.Println("--> Activity finished:", activity.ID)
				finishedCounter++
				cancel()

			default:
				if finishedCounter == startedCounter {
					return
				}
			}
		}
	}(finishCh)

	for startTime.Before(end) {
		activity := start(categories, startTime)

		semaphore <- struct{}{}

		go func(f chan<- *models.Activity, c int) {
			defer func() { <-semaphore }()

			timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
			defer cancel()

			started, err := actRepo.Create(timeoutCtx, activity)
			if err != nil {
				log.Panicln(err.Error())
			}

			log.Printf("%d - Activity started: %d\n", c, started.ID)
			f <- started

		}(finishCh, startedCounter)

		// incr
		startTime = startTime.Add(time.Hour)
		startedCounter++
	}

	log.Println("[ ] Total started: ", startedCounter)

	close(semaphore)

	wg.Wait()

	last := start(categories, time.Now())
	_, err = actRepo.Create(ctx, last)
	if err != nil {
		log.Panicln(err)
	}
}

func start(categories []*models.Category, t time.Time) *models.Activity {
	r := random().Intn(len(categories))
	category := categories[r]
	activity := &models.Activity{
		CategoryID:  category.ID,
		Description: fmt.Sprintf("IT_%s", uuid.NewString()),
		Status:      models.Started,
	}
	activity.SetStartedAt(t.Add(time.Hour * -1))
	activity.SetUpdatedAt(t.Add(time.Hour * -1))
	return activity
}

func random() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

type Counter struct {
	total uint32
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Inc() {
	atomic.AddUint32(&c.total, 1)
}

func (c *Counter) Total() int {
	total := atomic.LoadUint32(&c.total)
	return int(total)
}

func (c *Counter) Equals(other *Counter) bool {
	return c.total == other.total
}
