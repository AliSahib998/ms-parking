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
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strconv"
	"time"
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

	//for adding random parking slots
	loadParkingSlots(redisClient)

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

func loadParkingSlots(client *redis.Client) {
	var slots []model.Slot
	for i := 0; i < 100; i++ {
		var slotType = model.FourWheeler
		if i%3 == 0 {
			slotType = model.TwoWheeler
		}
		var slot = model.Slot{
			SlotNumber: rand.New(
				rand.NewSource(time.Now().UnixNano())).Int(),
			IsEmpty:  false,
			SlotType: slotType,
			Vehicle:  nil,
		}
		slots = append(slots, slot)
	}

	r, _ := json.Marshal(slots)
	client.Set(context.TODO(), model.SLOTS, r, 0)
}
