package http

import (
	d "app/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type UserHandler struct {
	su d.ServiceUser
}

// AddUser handles the request to add a new user
func (h *UserHandler) AddUser(c echo.Context) error {
	var u *d.User

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: http.StatusText(http.StatusBadRequest)})
	}

	// Validate the request body
	// Use the validator package to validate the struct fields
	validate = validator.New()
	if err := validate.Struct(u); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	_, err := h.su.AddUser(c.Request().Context(), u)
	if err != nil {
		if err == d.ErrConflictUsername {
			return c.JSON(http.StatusConflict, ResponseError{Message: err.Error()})
		} else {
			return c.JSON(http.StatusInternalServerError, ResponseError{Message: http.StatusText(http.StatusInternalServerError)})
		}
	}

	return c.JSON(http.StatusCreated, ResponseSuccess{Message: http.StatusText(http.StatusCreated)})
}

// Login handles the request to validate a user by username and password
func (h *UserHandler) Login(c echo.Context) error {
	var u *d.User

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: http.StatusText(http.StatusBadRequest)})
	}

	// Validate the request body
	// Use the validator package to validate the struct fields
	validate = validator.New()
	if err := validate.Struct(u); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	jwt, err := h.su.Login(c.Request().Context(), u)
	if err != nil {
		if err == d.ErrUserNotFound {
			return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
		} else if err == d.ErrLogin {
			return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
		} else {
			return c.JSON(http.StatusInternalServerError, ResponseError{Message: http.StatusText(http.StatusInternalServerError)})
		}
	}

	return c.JSON(http.StatusOK, ResponseSuccess{Token: jwt, Message: http.StatusText(http.StatusOK)})
}
