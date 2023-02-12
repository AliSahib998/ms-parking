package errhandler

type CommonError struct {
	Message string            `json:"message"`
	Code    string            `json:"code"`
	Errors  []ValidationError `json:"errors"`
}

type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type BadRequestError CommonError

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewBadRequestError(message string, errors []ValidationError) error {
	return &BadRequestError{
		Message: message,
		Code:    "invalid.request",
		Errors:  errors,
	}
}

type NotFoundError CommonError

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(message string, errors []ValidationError) error {
	return &BadRequestError{
		Message: message,
		Code:    "not_found",
		Errors:  errors,
	}
}

type PaymentError CommonError

func (e *PaymentError) Error() string {
	return e.Message
}

func NewPaymentError(message string, errors []ValidationError) error {
	return &BadRequestError{
		Message: message,
		Code:    "payment_error",
		Errors:  errors,
	}
}

type AlreadyParkedError CommonError

func (e *AlreadyParkedError) Error() string {
	return e.Message
}

func NewAlreadyParkedError(message string, errors []ValidationError) error {
	return &BadRequestError{
		Message: message,
		Code:    "already_parked",
		Errors:  errors,
	}
}
