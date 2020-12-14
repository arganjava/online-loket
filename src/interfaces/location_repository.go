package interfaces

import (
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/models"
)

type LocationRepository interface {
	CreateLocation(dto.LocationRequest) (int64, error)
	FindLocationById(id string) (*models.Location, error)
}
