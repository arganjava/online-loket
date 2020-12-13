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

	transactionService := NewTransactionService(db, eventService)
	transactionRouter := NewTransactionRouter(transactionService)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/location/create", locationRouter.CreateLocation)
		v1.POST("/event/create", eventRouter.CreateEvent)
		v1.POST("/event/ticket/create", eventRouter.CreateEventTicket)
		v1.GET("/event/get_info/:id", eventRouter.GetEventInfo)
		v1.POST("/transaction/purchase", transactionRouter.CreateTransaction)
		v1.GET("/transaction/get_info/:id", transactionRouter.GetTransactionInfo)
	}
	return r
}
