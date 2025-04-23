package http

import (
	d "app/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type LoanHandler struct {
	su d.ServiceLoan
}

// AddLoan handles the request to add a new loan
func (h *LoanHandler) AddLoan(c echo.Context) error {
	var l *d.Loan

	if err := c.Bind(&l); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: http.StatusText(http.StatusBadRequest)})
	}

	// Validate the request body
	// Use the validator package to validate the struct fields
	validate = validator.New()
	if err := validate.Struct(l); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	_, err := h.su.AddLoan(c.Request().Context(), l)
	if err != nil {
		if err == d.ErrLimitExceeded {
			return c.JSON(http.StatusConflict, ResponseError{Message: err.Error()})
		} else {
			return c.JSON(http.StatusInternalServerError, ResponseError{Message: http.StatusText(http.StatusInternalServerError)})
		}
	}

	return c.JSON(http.StatusCreated, ResponseSuccess{Message: http.StatusText(http.StatusCreated)})
}
