package service

import (
	"context"
	"github.com/AliSahib998/ms-parking/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCalculatePayment_Success(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = PaymentServiceBuilder(&redisHelper)
	var ticket = model.Ticket{
		SlotNumber:    267829,
		VehicleNumber: "AZ-100",
		TicketNumber:  "12334444",
		TicketStatus:  "ACTIVE",
		GivenDate:     time.Now(),
		VehicleSize:   "TWO_WHEELER",
		PaymentInfo: &model.PaymentInfo{
			TransactionId: "333434333",
			Price:         "0",
			PaymentStatus: "UNPAID",
			PaymentDate:   nil,
		},
	}
	redisHelper.On("GetTicket", context.TODO(), mock.Anything).Return(&ticket, nil)
	res, err := service.CalculatePayment(context.Background(), "23333")
	assert.NoError(t, err)
	assert.Equal(t, res.PaymentStatus, model.UnPaid)
	assert.Equal(t, res.TransactionId, "333434333")
}

func TestCalculatePayment_Error(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = PaymentServiceBuilder(&redisHelper)
	var ticket = model.Ticket{
		SlotNumber:    267829,
		VehicleNumber: "AZ-100",
		TicketNumber:  "12334444",
		TicketStatus:  "INACTIVE",
		GivenDate:     time.Now(),
		VehicleSize:   "TWO_WHEELER",
		PaymentInfo: &model.PaymentInfo{
			TransactionId: "333434333",
			Price:         "0",
			PaymentStatus: "UNPAID",
			PaymentDate:   nil,
		},
	}
	redisHelper.On("GetTicket", context.TODO(), mock.Anything).Return(&ticket, nil)
	_, err := service.CalculatePayment(context.Background(), "23333")
	assert.Equal(t, err.Error(), "ticket status is inactive or payment is already paid")
}

func TestMakePayment_Success(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = PaymentServiceBuilder(&redisHelper)
	var ticket = model.Ticket{
		SlotNumber:    267829,
		VehicleNumber: "AZ-100",
		TicketNumber:  "12334444",
		TicketStatus:  "ACTIVE",
		GivenDate:     time.Now(),
		VehicleSize:   "TWO_WHEELER",
		PaymentInfo: &model.PaymentInfo{
			TransactionId: "333434333",
			Price:         "0",
			PaymentStatus: "UNPAID",
			PaymentDate:   nil,
		},
	}

	var performPaymentRequest = model.PerformPaymentRequest{
		TicketNumber: "3332222",
		Price:        3.23,
	}
	redisHelper.On("GetTicket", context.TODO(), mock.Anything).Return(&ticket, nil)
	redisHelper.On("Set", context.TODO(), mock.Anything, mock.Anything, mock.Anything).Return(nil)
	err := service.MakePayment(context.Background(), performPaymentRequest)
	assert.NoError(t, err)
	assert.Equal(t, ticket.TicketStatus, model.Inactive)
	assert.Equal(t, ticket.PaymentInfo.PaymentStatus, model.Paid)
}

func TestMakePayment_AlreadyPaid(t *testing.T) {
	var redisHelper = MockRedisHelper{}
	var service = PaymentServiceBuilder(&redisHelper)
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
			Price:         "0",
			PaymentStatus: "PAID",
			PaymentDate:   &currentTime,
		},
	}

	var performPaymentRequest = model.PerformPaymentRequest{
		TicketNumber: "3332222",
		Price:        3.23,
	}
	redisHelper.On("GetTicket", context.TODO(), mock.Anything).Return(&ticket, nil)
	err := service.MakePayment(context.Background(), performPaymentRequest)
	assert.Equal(t, err.Error(), "ticket status is inactive or payment is already paid")
}
