package http

import (
	d "app/domain"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

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
	var lt d.LimitType
	if err := c.Bind(&lt); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: http.StatusText(http.StatusBadRequest)})
	}
	_, err := h.slt.AddLimitType(c.Request().Context(), &lt)
	fmt.Println(err)
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
