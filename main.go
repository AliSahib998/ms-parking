package main

import (
	"context"
	"encoding/json"
	"github.com/AliSahib998/ms-parking/cache"
	"github.com/AliSahib998/ms-parking/config"
	_ "github.com/AliSahib998/ms-parking/docs"
	"github.com/AliSahib998/ms-parking/handler"
	"github.com/AliSahib998/ms-parking/model"
	"github.com/AliSahib998/ms-parking/service"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// @title Api Documentation
// @version 1.0
// @host localhost:81
// @Accept json
// @Produce json
// @BasePath /
func main() {
	config.LoadConfig()
	router := chi.NewRouter()
	redisClient := cache.NewRedisClientForDev()
	var slot = []model.Slot{
		{
			SlotNumber: 267829,
			IsEmpty:    true,
			SlotType:   model.TwoWheeler,
			Vehicle:    nil,
		},

		{
			SlotNumber: 22222,
			IsEmpty:    true,
			SlotType:   model.TwoWheeler,
			Vehicle:    nil,
		},

		{
			SlotNumber: 222221,
			IsEmpty:    true,
			SlotType:   model.TwoWheeler,
			Vehicle:    nil,
		},

		{
			SlotNumber: 2224221,
			IsEmpty:    true,
			SlotType:   model.TwoWheeler,
			Vehicle:    nil,
		},
	}
	r, _ := json.Marshal(slot)
	redisClient.Set(context.TODO(), model.SLOTS, r, 0)

	var redisHelper = &service.RedisHelper{
		RedisClient: redisClient,
	}
	var parkingService = service.ParkingServiceBuilder(redisHelper)
	var paymentService = service.PaymentServiceBuilder(redisHelper)
	handler.NewParkingHandler(router, parkingService)
	handler.NewPaymentHandler(router, paymentService)
	handler.NewHealthHandler(router)
	port := strconv.Itoa(config.Props.Port)
	log.Info("Starting server at port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
