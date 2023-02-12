package errhandler

import (
	"github.com/AliSahib998/ms-parking/util"
	"net/http"
)

func ErrorHandler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute the final handler, and deal with errors
		err := h(w, r)
		if err != nil {
			switch e := err.(type) {
			case *BadRequestError:
				// We can retrieve the status here and write out a specific
				// HTTP status code.
				w.WriteHeader(http.StatusBadRequest)
				util.Encode(w, e)
			case *NotFoundError:
				// We can retrieve the status here and write out a specific
				// HTTP status code.
				w.WriteHeader(http.StatusNotFound)
				util.Encode(w, e)
			case *PaymentError:
				// We can retrieve the status here and write out a specific
				// HTTP status code.
				w.WriteHeader(http.StatusBadRequest)
				util.Encode(w, e)
			default:
				// Any error types we don't specifically look out for default
				// to serving an HTTP 500
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}
	}
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error
