package repository

import (
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/models"
	mock "github.com/stretchr/testify/mock"
)

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) CreateEvent(request dto.EventRequest) (int64, error) {
	call := m.Called(request)
	res := call.Get(0)
	if res == nil {
		return 0, call.Error(1)
	}
	return 1, nil
}

func (m *MockEventRepository) FindEventById(id string) (*models.Event, error) {
	call := m.Called(id)
	res := call.Get(0)
	if res == nil {
		return nil, call.Error(1)
	}
	return res.(*models.Event), call.Error(1)
}

func (m *MockEventRepository) FindEventTicketId(id string) (*models.EventTicket, error) {
	call := m.Called(id)
	res := call.Get(0)
	if res == nil {
		return nil, call.Error(1)
	}
	return res.(*models.EventTicket), call.Error(1)
}

func (m *MockEventRepository) CreateEventTicket(request dto.EventTicketRequest) (int64, error) {
	call := m.Called(request)
	res := call.Get(0)
	if res == nil {
		return 0, call.Error(1)
	}
	return 1, nil
}
