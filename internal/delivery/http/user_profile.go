package http

import (
	d "app/domain"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type UserProfileHandler struct {
	sup d.ServiceUserProfile
}

// AddUserProfile handles the request to add a user_profile
func (h *UserProfileHandler) AddUserProfile(c echo.Context) error {
	var up *d.UserProfile

	if err := c.Bind(&up); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: http.StatusText(http.StatusBadRequest)})
	}

	up.UserID = int(c.Get("userID").(int64))

	// Validate the request body
	// Use the validator package to validate the struct fields
	validate = validator.New()
	if err := validate.Struct(up); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	_, err := h.sup.AddUserProfile(c.Request().Context(), up)
	fmt.Println(err)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: http.StatusText(http.StatusInternalServerError)})
	}

	return c.JSON(http.StatusCreated, ResponseSuccess{Message: http.StatusText(http.StatusCreated)})
}
