package http

import (
	d "app/domain"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

var validate *validator.Validate

type ResponseSuccess struct {
	Data    interface{} `json:"data,omitempty"`
	Token   string      `json:"token,omitempty"`
	Message string      `json:"message"`
}

type ResponseError struct {
	Message string `json:"message"`
}

// registers LimitType routes with the provided Echo instance.
func NewLimitTypeHandler(e *echo.Echo, slt d.ServiceLimitType) {
	handler := &LimitTypeHandler{slt}
	//adding
	e.POST("/limittype", handler.AddLimitType)
	//listing
	e.GET("/limittypes", handler.GetLimitTypes)
}

// registers User routes with the provided Echo instance.
func NewUserHandler(e *echo.Echo, su d.ServiceUser) {
	handler := &UserHandler{su}
	//adding
	e.POST("/register", handler.AddUser)

	//validating
	e.POST("/login", handler.Login)
}
