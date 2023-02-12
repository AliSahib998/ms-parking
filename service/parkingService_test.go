package service

import (
	"context"
	"errors"
	"github.com/AliSahib998/ms-parking/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

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
		SlotType:   model.FourWheeler,
		Vehicle:    nil,
	},
}

func TestGetParkingSlots_Success(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = ParkingServiceBuilder(&redisHelper)
	redisHelper.On("GetSlots", context.TODO(), model.SLOTS).Return(slot, nil)
	res, err := service.GetParkingSlots(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, len(res), 2)
}

func TestGetParkingSlots_Error(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = ParkingServiceBuilder(&redisHelper)
	redisHelper.On("GetSlots", context.TODO(), model.SLOTS).Return([]model.Slot{}, errors.New("slots are not found"))
	res, err := service.GetParkingSlots(context.Background())
	assert.Error(t, err, "slots are not found")
	assert.Equal(t, len(res), 0)
}

func TestGetParkingNumber_Success(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = ParkingServiceBuilder(&redisHelper)
	redisHelper.On("GetSlots", context.TODO(), model.SLOTS).Return(slot, nil)
	redisHelper.On("Set", context.TODO(), model.SLOTS, mock.Anything, mock.Anything).Return(nil)
	redisHelper.On("Set", context.TODO(), mock.Anything, mock.Anything, mock.Anything).Return(nil)
	var vehicle = model.Vehicle{
		VehicleNumber: "AZ-100",
		VehicleSize:   "TWO_WHEELER",
	}
	res, err := service.GetParkingNumber(context.Background(), &vehicle)
	assert.NoError(t, err)
	assert.Equal(t, res.VehicleNumber, vehicle.VehicleNumber)
}

func TestGetParkingNumber_AlreadyParked(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = ParkingServiceBuilder(&redisHelper)
	var slot = []model.Slot{
		{
			SlotNumber: 267829,
			IsEmpty:    false,
			SlotType:   model.TwoWheeler,
			Vehicle: &model.Vehicle{
				VehicleNumber: "AZ-100",
				VehicleSize:   "TWO_WHEELER",
			},
		},
	}
	redisHelper.On("GetSlots", context.TODO(), model.SLOTS).Return(slot, nil)
	var vehicle = model.Vehicle{
		VehicleNumber: "AZ-100",
		VehicleSize:   "TWO_WHEELER",
	}
	_, err := service.GetParkingNumber(context.Background(), &vehicle)
	assert.Equal(t, err.Error(), "vehicle AZ-100 is already parked")
}

func TestLeaveParkingSlot_Success(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = ParkingServiceBuilder(&redisHelper)
	var currentTime = time.Now()
	var ticket = model.Ticket{
		SlotNumber:    267829,
		VehicleNumber: "AZ-100",
		TicketNumber:  "12334444",
		TicketStatus:  "ACTIVE",
		GivenDate:     time.Now(),
		VehicleSize:   "TWO_WHEELER",
		PaymentInfo: &model.PaymentInfo{
			TransactionId: "333434333",
			Price:         "3.43",
			PaymentStatus: "PAID",
			PaymentDate:   &currentTime,
		},
	}
	redisHelper.On("GetTicket", context.TODO(), mock.Anything).Return(&ticket, nil)
	redisHelper.On("GetSlots", context.TODO(), model.SLOTS).Return(slot, nil)
	redisHelper.On("Set", context.TODO(), model.SLOTS, mock.Anything, mock.Anything).Return(nil)
	err := service.LeaveParkingSlot(context.Background(), "2344545")
	assert.NoError(t, err)
}

func TestLeaveParkingSlot_NeededPayment(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = ParkingServiceBuilder(&redisHelper)
	var ticket = model.Ticket{
		SlotNumber:    267829,
		VehicleNumber: "AZ-100",
		TicketNumber:  "12334444",
		TicketStatus:  "ACTIVE",
		GivenDate:     time.Now(),
		VehicleSize:   "TWO_WHEELER",
		PaymentInfo: &model.PaymentInfo{
			TransactionId: "333434333",
			Price:         "3.43",
			PaymentStatus: "UNPAID",
			PaymentDate:   nil,
		},
	}
	redisHelper.On("GetTicket", context.TODO(), mock.Anything).Return(&ticket, nil)
	redisHelper.On("GetSlots", context.TODO(), model.SLOTS).Return(slot, nil)
	redisHelper.On("Set", context.TODO(), model.SLOTS, mock.Anything, mock.Anything).Return(nil)
	err := service.LeaveParkingSlot(context.Background(), "2344545")
	assert.Equal(t, err.Error(), "needed to payment")
}
