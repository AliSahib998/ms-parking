package service

import (
	"context"
	"encoding/json"
	"fmt"
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
	redisHelper IRedisHelper
}

func ParkingServiceBuilder(redisHelper IRedisHelper) *ParkingService {
	return &ParkingService{
		redisHelper: redisHelper,
	}
}

func (p *ParkingService) GetParkingSlots(ctx context.Context) ([]model.Slot, error) {
	return p.redisHelper.GetSlots(ctx, model.SLOTS)
}

func (p *ParkingService) GetParkingNumber(ctx context.Context, vehicle *model.Vehicle) (*model.Ticket, error) {
	var slots, err = p.redisHelper.GetSlots(ctx, model.SLOTS)
	if err != nil {
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

	//check ticket is null or no
	if ticket == nil {
		log.Info("not found any empty slot")
		return nil, errhandler.NewNotFoundError("any empty slot not found", nil)
	}

	//update slots and add ticket
	slotByteArray, _ := json.Marshal(slots)
	ticketByteArray, _ := json.Marshal(ticket)
	_ = p.redisHelper.Set(ctx, model.SLOTS, slotByteArray, 0)
	_ = p.redisHelper.Set(ctx, ticket.TicketNumber, ticketByteArray, 0)
	return ticket, nil
}

func (p *ParkingService) LeaveParkingSlot(ctx context.Context, ticketNumber string) error {
	//get ticket from database
	var ticket, err = p.redisHelper.GetTicket(ctx, ticketNumber)
	if err != nil {
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
	slots, err := p.redisHelper.GetSlots(ctx, model.SLOTS)

	for i, _ := range slots {
		if slots[i].SlotNumber == slotNumber {
			slots[i].IsEmpty = true
			slots[i].Vehicle = nil
			break
		}
	}
	slotByteArray, _ := json.Marshal(slots)
	err = p.redisHelper.Set(ctx, model.SLOTS, slotByteArray, 0)
	return err
}
