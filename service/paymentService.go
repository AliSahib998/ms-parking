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

type IPaymentService interface {
	CalculatePayment(ctx context.Context, ticketNumber string) (*model.PaymentInfo, error)
	MakePayment(ctx context.Context, paymentRequest model.PerformPaymentRequest) error
}

type PaymentService struct {
	redisClient cache.IRedisClient
}

func PaymentServiceBuilder(redisClient cache.IRedisClient) *PaymentService {
	return &PaymentService{
		redisClient: redisClient,
	}
}

func (p *PaymentService) CalculatePayment(ctx context.Context, ticketNumber string) (*model.PaymentInfo, error) {
	//get ticket from database
	var ticket *model.Ticket
	var value, err = p.redisClient.Get(ctx, ticketNumber).Result()
	if len(value) == 0 {
		log.Info("this key is not exist in redis %s", ticketNumber)
		return nil, errhandler.NewNotFoundError(fmt.Sprintf("ticket number %s is not found", ticketNumber), nil)
	}
	err = json.Unmarshal([]byte(value), &ticket)
	if err != nil {
		log.Error("error occurred when unmarshalling:", err)
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
	//get ticket from database
	var ticket *model.Ticket
	var value, err = p.redisClient.Get(ctx, paymentRequest.TicketNumber).Result()
	if len(value) == 0 {
		log.Info("this key is not exist in redis %s", paymentRequest.TicketNumber)
		return errhandler.NewNotFoundError(fmt.Sprintf("ticket number %s is not found",
			paymentRequest.TicketNumber), nil)
	}
	err = json.Unmarshal([]byte(value), &ticket)
	if err != nil {
		log.Error("error occurred when unmarshalling:", err)
		return err
	}

	//check ticket status
	if ticket.TicketStatus == model.Inactive || ticket.PaymentInfo.PaymentStatus == model.Paid {
		return errhandler.NewPaymentError("ticket status is inactive or payment is already paid", nil)
	}

	//perform payment (call any payment service)
	ticket.PaymentInfo.Price = fmt.Sprintf("%f", paymentRequest.Price)
	ticket.PaymentInfo.PaymentStatus = model.Paid
	ticket.TicketStatus = model.Inactive

	//update ticket in redis
	ticketByteArray, _ := json.Marshal(ticket)
	p.redisClient.Set(ctx, ticket.TicketNumber, ticketByteArray, 0)
	return nil
}

func calculateHoursPrice(seconds float64) string {
	if seconds <= 3600 {
		return "0"
	}
	return fmt.Sprintf("%.2f", ((seconds-3600)*2)/3600)
}
