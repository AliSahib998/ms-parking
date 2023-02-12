package handler

import (
	"encoding/json"
	"fmt"
	"github.com/AliSahib998/ms-parking/errhandler"
	"github.com/AliSahib998/ms-parking/model"
	"github.com/AliSahib998/ms-parking/service"
	"github.com/AliSahib998/ms-parking/util"
	"github.com/go-chi/chi"
	"net/http"
)

type paymentHandler struct {
	paymentService service.IPaymentService
}

func NewPaymentHandler(router *chi.Mux, paymentService service.IPaymentService) *chi.Mux {
	h := &paymentHandler{paymentService: paymentService}
	router.Group(func(r chi.Router) {
		//router.Use(middleware.RequestParamsMiddleware)
		router.Get("/payment/calculate/{ticketNumber}", errhandler.ErrorHandler(h.CalculatePayment))
		router.Post("/payment/perform", errhandler.ErrorHandler(h.PerformPayment))
	})
	return router
}

func (h *paymentHandler) CalculatePayment(w http.ResponseWriter, r *http.Request) error {
	ticketNumber := chi.URLParam(r, "ticketNumber")
	if len(ticketNumber) == 0 {
		return errhandler.NewBadRequestError(fmt.Sprintf("%s", "invalid ticket number"), nil)
	}
	paymentInfo, err := h.paymentService.CalculatePayment(r.Context(), ticketNumber)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(paymentInfo)
	return err
}

func (h *paymentHandler) PerformPayment(w http.ResponseWriter, r *http.Request) error {
	request := new(model.PerformPaymentRequest)
	err := util.ParseRequest(r, request)
	if err != nil {
		return err
	}

	err = h.paymentService.MakePayment(r.Context(), *request)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	return nil
}
