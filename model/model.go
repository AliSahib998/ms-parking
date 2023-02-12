package model

import "time"

type VehicleSize string

type PaymentStatus string

type TicketStatus string

const (
	TwoWheeler  VehicleSize   = "TWO_WHEELER"
	FourWheeler VehicleSize   = "FOUR_WHEELER"
	Paid        PaymentStatus = "PAID"
	UnPaid      PaymentStatus = "UNPAID"
	Active      TicketStatus  = "ACTIVE"
	Inactive    TicketStatus  = "INACTIVE"
)

type Vehicle struct {
	VehicleNumber string      `json:"vehicleNumber"`
	VehicleSize   VehicleSize `json:"vehicleSize"`
}

type Slot struct {
	SlotNumber int         `json:"slotNumber"`
	IsEmpty    bool        `json:"IsEmpty"`
	SlotType   VehicleSize `json:"slotType"`
	Vehicle    *Vehicle    `json:"vehicle"`
}

type Ticket struct {
	SlotNumber    int          `json:"slotNumber"`
	VehicleNumber string       `json:"vehicleNumber"`
	TicketNumber  string       `json:"ticketNumber"`
	TicketStatus  TicketStatus `json:"ticketStatus"`
	GivenDate     time.Time    `json:"givenDate"`
	VehicleSize   VehicleSize  `json:"vehicleSize"`
	PaymentInfo   *PaymentInfo `json:"paymentInfo"`
}

type PaymentInfo struct {
	TransactionId string        `json:"transactionId"`
	Price         string        `json:"price"`
	PaymentStatus PaymentStatus `json:"paymentStatus"`
	PaymentDate   *time.Time    `json:"paymentDate"`
}

type PerformPaymentRequest struct {
	TicketNumber string  `json:"ticketNumber"`
	Price        float64 `json:"price"`
}
