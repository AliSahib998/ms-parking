package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AliSahib998/ms-parking/cache"
	"github.com/AliSahib998/ms-parking/errhandler"
	"github.com/AliSahib998/ms-parking/model"
	"github.com/AliSahib998/ms-parking/util"
	log "github.com/sirupsen/logrus"
	"time"
)

type IParkingService interface {
	GetParkingSlots(ctx context.Context) ([]model.Slot, error)
	GetParkingNumber(ctx context.Context, vehicle *model.Vehicle) (*model.Ticket, error)
	LeaveParkingSlot(ctx context.Context, ticketNumber string) error
}

type ParkingService struct {
	redisClient cache.IRedisClient
}

func ParkingBuilder(redisClient cache.IRedisClient) *ParkingService {
	return &ParkingService{
		redisClient: redisClient,
	}
}

func (p *ParkingService) GetParkingSlots(ctx context.Context) ([]model.Slot, error) {
	var slots []model.Slot
	var value, err = p.redisClient.Get(ctx, model.SLOTS).Result()
	if len(value) == 0 {
		log.Info("this key is not exist in redis %s", model.SLOTS)
		return slots, errhandler.NewNotFoundError("slots is not found", nil)
	}
	err = json.Unmarshal([]byte(value), &slots)
	if err != nil {
		log.Error("error occurred when unmarshalling:", err)
	}

	return slots, err
}

func (p *ParkingService) GetParkingNumber(ctx context.Context, vehicle *model.Vehicle) (*model.Ticket, error) {
	var slots []*model.Slot
	var value, err = p.redisClient.Get(ctx, model.SLOTS).Result()
	if len(value) == 0 {
		log.Info("this key is not exist in redis %s", model.SLOTS)
		return nil, errhandler.NewNotFoundError("slots is not found", nil)
	}
	err = json.Unmarshal([]byte(value), &slots)
	if err != nil {
		log.Error("error occurred when unmarshalling:", err)
		return nil, err
	}

	var ticket *model.Ticket

	//check vehicle number
	for _, v := range slots {
		if v.Vehicle != nil && v.Vehicle.VehicleNumber == vehicle.VehicleNumber {
			return nil, errhandler.NewAlreadyParkedError(fmt.Sprintf("vehicle %s is already parked", vehicle.VehicleNumber), nil)
		}
	}

	//find empty slot
	for i, v := range slots {
		if v.IsEmpty && v.SlotType == vehicle.VehicleSize {
			slots[i].IsEmpty = false
			slots[i].Vehicle = vehicle
			ticket = &model.Ticket{
				SlotNumber:    v.SlotNumber,
				VehicleNumber: vehicle.VehicleNumber,
				TicketNumber:  util.RandomNumber(),
				GivenDate:     time.Now(),
				TicketStatus:  model.Active,
				VehicleSize:   vehicle.VehicleSize,
				PaymentInfo: &model.PaymentInfo{
					TransactionId: util.RandomNumber(),
					PaymentStatus: model.UnPaid,
					Price:         "0",
					PaymentDate:   nil,
				},
			}
			break
		}
	}

	if ticket == nil {
		log.Info("not found any empty slot")
		return nil, errhandler.NewNotFoundError("any empty slot not found", nil)
	}

	//update
	slotByteArray, _ := json.Marshal(slots)
	ticketByteArray, _ := json.Marshal(ticket)
	p.redisClient.Set(ctx, model.SLOTS, slotByteArray, 0)
	p.redisClient.Set(ctx, ticket.TicketNumber, ticketByteArray, 0)
	return ticket, nil
}

func (p *ParkingService) LeaveParkingSlot(ctx context.Context, ticketNumber string) error {
	//get ticket from database
	var ticket *model.Ticket
	var value, err = p.redisClient.Get(ctx, ticketNumber).Result()
	if len(value) == 0 {
		log.Info("this key is not exist in redis %s", ticketNumber)
		return errhandler.NewNotFoundError(fmt.Sprintf("ticket number %s is not found", ticketNumber), nil)
	}
	err = json.Unmarshal([]byte(value), &ticket)
	if err != nil {
		log.Error("error occurred when unmarshalling:", err)
		return err
	}

	//check payment status for car leaving
	if (ticket.PaymentInfo.PaymentStatus == model.UnPaid) ||
		(ticket.PaymentInfo.PaymentStatus == model.Paid &&
			time.Now().Sub(*ticket.PaymentInfo.PaymentDate).Minutes() >= 10) {
		return errhandler.NewPaymentError("needed to payment", nil)
	}

	//update slot status in parking
	var slotNumber = ticket.SlotNumber
	var slots []*model.Slot
	var slotValue, slotErr = p.redisClient.Get(ctx, model.SLOTS).Result()
	if len(slotValue) == 0 {
		log.Info("this key is not exist in redis %s", model.SLOTS)
		return errhandler.NewNotFoundError("slots were not found", nil)
	}
	slotErr = json.Unmarshal([]byte(value), &slots)
	if slotErr != nil {
		log.Error("error occurred when unmarshalling:", err)
		return err
	}

	for i, _ := range slots {
		if slots[i].SlotNumber == slotNumber {
			slots[i].IsEmpty = true
			slots[i].Vehicle = nil
			break
		}
	}

	//update slots in redis
	slotByteArray, _ := json.Marshal(slots)
	p.redisClient.Set(ctx, model.SLOTS, slotByteArray, 0)
	return nil
}
