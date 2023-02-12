package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AliSahib998/ms-parking/errhandler"
	"github.com/AliSahib998/ms-parking/model"
	"time"
)

type IPaymentService interface {
	CalculatePayment(ctx context.Context, ticketNumber string) (*model.PaymentInfo, error)
	MakePayment(ctx context.Context, paymentRequest model.PerformPaymentRequest) error
}

type PaymentService struct {
	redisHelper IRedisHelper
}

func PaymentServiceBuilder(redisHelper IRedisHelper) *PaymentService {
	return &PaymentService{
		redisHelper: redisHelper,
	}
}

func (p *PaymentService) CalculatePayment(ctx context.Context, ticketNumber string) (*model.PaymentInfo, error) {
	//get ticket from database
	var ticket, err = p.redisHelper.GetTicket(ctx, ticketNumber)
	if err != nil {
		return nil, err
	}
	//check ticket status
	if ticket.TicketStatus == model.Inactive || ticket.PaymentInfo.PaymentStatus == model.Paid {
		return nil, errhandler.NewPaymentError("ticket status is inactive or payment is already paid", nil)
	}

	//calculate hours and price
	var parkingTime = time.Now().Sub(ticket.GivenDate).Seconds()
	var price = calculateHoursPrice(parkingTime)
	var paymentInfo = ticket.PaymentInfo
	paymentInfo.Price = price
	return paymentInfo, nil
}

func (p *PaymentService) MakePayment(ctx context.Context, paymentRequest model.PerformPaymentRequest) error {
	//get ticket from database
	var ticket, err = p.redisHelper.GetTicket(ctx, paymentRequest.TicketNumber)
	if err != nil {
		return err
	}

	//check ticket status
	if ticket.TicketStatus == model.Inactive || ticket.PaymentInfo.PaymentStatus == model.Paid {
		return errhandler.NewPaymentError("ticket status is inactive or payment is already paid", nil)
	}

	//perform payment (call any payment service)
	ticket.PaymentInfo.Price = fmt.Sprintf("%f", paymentRequest.Price)
	ticket.PaymentInfo.PaymentStatus = model.Paid
	var currentTime = time.Now()
	ticket.PaymentInfo.PaymentDate = &currentTime
	ticket.TicketStatus = model.Inactive

	//update ticket in redis
	ticketByteArray, _ := json.Marshal(ticket)
	err = p.redisHelper.Set(ctx, ticket.TicketNumber, ticketByteArray, 0)
	return nil
}

func calculateHoursPrice(seconds float64) string {
	//10 minutes = 3600 seconds
	if seconds <= 3600 {
		return "0"
	}
	return fmt.Sprintf("%.2f", ((seconds-3600)*2)/3600)
}
