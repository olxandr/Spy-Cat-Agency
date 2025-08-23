package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var InternalServerError = errors.New("Internal server error")

func writeJSONValidationErrors(c *gin.Context, errs map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
}
