package http

import (
	d "app/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type LimitTypeHandler struct {
	slt d.ServiceLimitType
}

// AddLimitType handles the request to add a new limit type
func (h *LimitTypeHandler) AddLimitType(c echo.Context) error {
	var lt *d.LimitType

	if err := c.Bind(&lt); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: http.StatusText(http.StatusBadRequest)})
	}

	// Validate the request body
	// Use the validator package to validate the struct fields
	validate = validator.New()
	if err := validate.Struct(lt); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	_, err := h.slt.AddLimitType(c.Request().Context(), lt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: http.StatusText(http.StatusInternalServerError)})
	}

	return c.JSON(http.StatusCreated, ResponseError{Message: http.StatusText(http.StatusCreated)})
}

// GetLimitTypes handles the request to get all limit types
func (h *LimitTypeHandler) GetLimitTypes(c echo.Context) error {
	lts, err := h.slt.GetLimitTypes(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: http.StatusText(http.StatusInternalServerError)})
	}

	return c.JSON(http.StatusOK, ResponseSuccess{Data: lts, Message: http.StatusText(http.StatusOK)})
}
