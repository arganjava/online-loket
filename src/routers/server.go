package routers

import (
	"database/sql"
	. "github.com/arganjava/online-loket/src/service"
	"github.com/gin-gonic/gin"
)

func SetupServer(db *sql.DB) *gin.Engine {
	locationService := NewLocationService(db)
	locationRouter := NewLocationRouter(locationService)

	eventService := NewEventService(db, locationService)
	eventRouter := NewEventRouter(eventService)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/location/create", locationRouter.CreateLocation)
		v1.POST("/event/create", eventRouter.CreateEvent)
	}
	return r
}
