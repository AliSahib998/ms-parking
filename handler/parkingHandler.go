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

type parkingHandler struct {
	parkingService service.IParkingService
}

func NewParkingHandler(router *chi.Mux, parkingService service.IParkingService) *chi.Mux {
	h := &parkingHandler{parkingService: parkingService}
	router.Group(func(r chi.Router) {
		//we can use this middleware for logging and adding header parameters
		//router.Use(middleware.RequestParamsMiddleware)
		router.Get("/parking/stat", errhandler.ErrorHandler(h.GetParkingSlots))
		router.Post("/parking/entrance", errhandler.ErrorHandler(h.GetParkingNumber))
		router.Get("/parking/depart/{ticketNumber}", errhandler.ErrorHandler(h.LeaveParkingSlot))
	})
	return router
}

func (h *parkingHandler) GetParkingSlots(w http.ResponseWriter, r *http.Request) error {
	slots, err := h.parkingService.GetParkingSlots(r.Context())
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(slots)
	return err
}

func (h *parkingHandler) GetParkingNumber(w http.ResponseWriter, r *http.Request) error {
	request := new(model.Vehicle)
	err := util.ParseRequest(r, request)
	if err != nil {
		return err
	}

	ticket, err := h.parkingService.GetParkingNumber(r.Context(), request)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ticket)
	return err
}

func (h *parkingHandler) LeaveParkingSlot(w http.ResponseWriter, r *http.Request) error {
	ticketNumber := chi.URLParam(r, "ticketNumber")
	if len(ticketNumber) == 0 {
		return errhandler.NewBadRequestError(fmt.Sprintf("%s", "invalid ticket number"), nil)
	}
	err := h.parkingService.LeaveParkingSlot(r.Context(), ticketNumber)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	return nil
}
