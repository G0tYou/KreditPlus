package http

import (
	d "app/domain"
	"app/internal/delivery/http/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
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

// registers UserProfile routes with the provided Echo instance.
func NewUserProfileHandler(e *echo.Echo, sup d.ServiceUserProfile) {
	handler := &UserProfileHandler{sup}

	userGroup := e.Group("/user")
	userGroup.Use(middleware.JWTAuth(viper.GetString("jwt.secret")))

	//adding
	userGroup.POST("/profile", handler.AddUserProfile)

}
