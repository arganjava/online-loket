package routers

import (
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/service_repository"
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
	"net/http"
)

type LocationRouter struct {
	locationService service_repository.LocationRepository
}

func NewLocationRouter(locationService service_repository.LocationRepository) *LocationRouter {
	return &LocationRouter{
		locationService: locationService,
	}
}

func (l LocationRouter) CreateLocation(c *gin.Context) {
	var location dto.LocationRequest
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := validator.Validate(location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	_, err := l.locationService.CreateLocation(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Location created successfully!"})

}
