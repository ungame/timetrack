package app

import (
	"context"
	"flag"
	"github.com/ungame/timetrack/app/handlers"
	"github.com/ungame/timetrack/app/observer"
	"github.com/ungame/timetrack/app/repository"
	"github.com/ungame/timetrack/app/router"
	"github.com/ungame/timetrack/app/service"
	"github.com/ungame/timetrack/db"
	"github.com/ungame/timetrack/httpext"
	"github.com/ungame/timetrack/ioext"
	"log"
	"net/http"
	"time"
)

var port int

func init() {
	flag.IntVar(&port, "p", 15555, "set port")
	flag.Parse()
}

func Run() {
	conn := db.New()
	defer ioext.Close(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		log.Panicln("unable to ping database:", err.Error())
	}

	var (
		categoriesRepository = repository.NewCategoriesRepository(conn)
		categoriesService    = service.NewCategoriesService(categoriesRepository)
		categoriesHandler    = handlers.NewCategoriesHandler(categoriesService)
		activitiesRepository = repository.NewActivitiesRepository(conn)
		activitiesObserver   = observer.New("activities")
		activitiesService    = service.NewActivitiesService(categoriesService, activitiesRepository, activitiesObserver)
		activitiesHandler    = handlers.NewActivitiesHandler(activitiesService)
	)

	r := router.New(categoriesHandler, activitiesHandler)

	log.Printf("Listening http://localhost:%d\n\n", port)
	log.Fatalln(http.ListenAndServe(httpext.Port(port).Addr(), r.WithCors()))
}
