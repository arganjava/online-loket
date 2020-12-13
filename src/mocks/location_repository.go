package repository

import (
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/models"
	mock "github.com/stretchr/testify/mock"
)

type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) CreateLocation(request dto.LocationRequest) (int64, error) {
	call := m.Called(request)
	res := call.Get(0)
	if res == nil {
		return 0, call.Error(1)
	}
	return 1, nil
}

// GetListTeller :
func (m *MockLocationRepository) FindLocationById(id string) (*models.Location, error) {
	call := m.Called(id)
	res := call.Get(0)
	if res == nil {
		return nil, call.Error(1)
	}
	return res.(*models.Location), call.Error(1)
}
