package http

import (
	d "app/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

var validate *validator.Validate

type ResponseError struct {
	Message string `json:"message"`
}

type LimitTypeHandler struct {
	slt d.ServiceLimitType
}

func NewLimitTypeHandler(e *echo.Echo, slt d.ServiceLimitType) {
	handler := &LimitTypeHandler{slt}
	//adding
	e.POST("/limittype", handler.AddLimitType)
	//listing
	e.GET("/limittypes", handler.GetLimitTypes)
}

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
		return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusCreated, http.StatusText(http.StatusCreated))
}

func (h *LimitTypeHandler) GetLimitTypes(c echo.Context) error {
	lts, err := h.slt.GetLimitTypes(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: http.StatusText(http.StatusInternalServerError)})
	}

	return c.JSON(http.StatusOK, lts)
}
