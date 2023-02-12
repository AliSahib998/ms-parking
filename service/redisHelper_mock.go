package service

import (
	"context"
	"github.com/AliSahib998/ms-parking/model"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRedisHelper struct {
	mock.Mock
}

func (m *MockRedisHelper) GetSlots(ctx context.Context, key string) ([]model.Slot, error) {
	args := m.Called(ctx, key)
	return args.Get(0).([]model.Slot), args.Error(1)
}

func (m *MockRedisHelper) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, val, expiration)
	return args.Error(0)
}

func (m *MockRedisHelper) GetTicket(ctx context.Context, key string) (*model.Ticket, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*model.Ticket), args.Error(1)
}
