package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AliSahib998/ms-parking/cache"
	"github.com/AliSahib998/ms-parking/errhandler"
	"github.com/AliSahib998/ms-parking/model"
	log "github.com/sirupsen/logrus"
	"time"
)

type IRedisHelper interface {
	GetSlots(ctx context.Context, key string) ([]model.Slot, error)
	Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error
	GetTicket(ctx context.Context, key string) (*model.Ticket, error)
}

type RedisHelper struct {
	RedisClient cache.IRedisClient
}

func (r *RedisHelper) GetSlots(ctx context.Context, key string) ([]model.Slot, error) {
	var slots []model.Slot
	var value, err = r.RedisClient.Get(ctx, key).Result()
	if len(value) == 0 {
		log.Infof("this key is not exist in redis %s", model.SLOTS)
		return slots, errhandler.NewNotFoundError("slots is not found", nil)
	}
	err = json.Unmarshal([]byte(value), &slots)
	if err != nil {
		log.Error("error occurred when unmarshalling:", err)
	}

	return slots, err
}

func (r *RedisHelper) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	var status = r.RedisClient.Set(ctx, key, val, expiration)
	return status.Err()
}

func (r *RedisHelper) GetTicket(ctx context.Context, key string) (*model.Ticket, error) {
	var ticket *model.Ticket
	var value, err = r.RedisClient.Get(ctx, key).Result()
	if len(value) == 0 {
		log.Infof("this key is not exist in redis %s", key)
		return nil, errhandler.NewNotFoundError(fmt.Sprintf("ticket number %s is not found", key), nil)
	}
	err = json.Unmarshal([]byte(value), &ticket)
	if err != nil {
		log.Error("error occurred when unmarshalling:", err)
		return nil, err
	}
	return ticket, nil
}
