package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var InternalServerError = errors.New("Internal server error")

func returnValidationErrors(c *gin.Context, errs validator.ValidationErrors) {
	errorMap := make(map[string]string)
	for _, e := range errs {
		errorMap[e.Field()] = e.Tag()
	}
	c.JSON(http.StatusUnprocessableEntity, errorMap)
}
