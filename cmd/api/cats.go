package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"spy-cat-agency/internal/models"
	"spy-cat-agency/internal/validator"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new cat
// @Description Create a new spy cat
// @Tags cats
// @Accept  json
// @Produce  json
// @Param cat body models.Cat true "Cat object"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cats/create [post]
func (app *application) createCat(c *gin.Context) {
	var (
		v   = validator.New()
		cat = &models.Cat{}
	)
	if err := c.ShouldBindJSON(cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	v.Check(cat.Name != "", "name", validator.ErrEmptyFIeld.Error())
	v.Check(cat.Breed != "", "breed", "can't be empty")
	v.Check(app.cats.Breeds.Exists(cat.Breed), "breed", "invalid breed")

	if !v.Valid() {
		writeJSONValidationErrors(c, v.Errors)
		return
	}
	id, err := app.cats.Create(c, cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error while creating a cat, try again later"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"info": "success", "id": id})

	slog.Info("Cat created", "id", id)
}

// @Summary Remove a cat
// @Description Remove a spy cat by ID
// @Tags cats
// @Accept  json
// @Produce  json
// @Param id path int true "Cat ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cats/remove/{id} [delete]
func (app *application) removeCat(c *gin.Context) {
	v := validator.New()
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	v.Check(id != 0, "id", validator.ErrZeroID.Error())

	if !v.Valid() {
		writeJSONValidationErrors(c, v.Errors)
		return
	}

	if err := app.cats.Remove(c, id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cat with ID %d doesn't exist", id)})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"info": "success"})

	slog.Info("Cat removed", "id", id)
}

// @Summary Update a cat's salary
// @Description Update a spy cat's salary by ID
// @Tags cats
// @Accept  json
// @Produce  json
// @Param cat body models.Cat true "Cat object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cats/update_salary [put]
func (app *application) updateCatsSalary(c *gin.Context) {
	cat := &models.Cat{}
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	v := validator.New()
	v.Check(cat.Salary != 0, "salary", validator.ErrEmptyFIeld.Error())
	v.Check(cat.ID != 0, "id", validator.ErrZeroID.Error())
	if !v.Valid() {
		writeJSONValidationErrors(c, v.Errors)
		return
	}

	updatedCat, err := app.cats.UpdateSalary(c, cat)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": "Cat doesn't exist"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"info": "succes", "updated cat": updatedCat})

	slog.Info("Cat's salary updated", "id", cat.ID, "before", cat.Salary, "after", updatedCat.Salary)
}

// @Summary List all cats
// @Description Get a list of all spy cats
// @Tags cats
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Cat
// @Failure 500 {object} map[string]interface{}
// @Router /cats/list [get]
func (app *application) listCats(c *gin.Context) {
	cats, err := app.cats.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error while fetching cats, try again later"})
		return
	}

	c.JSON(http.StatusOK, cats)

	slog.Info("Cats listed", "cats", len(cats))
}

// @Summary Get a cat by ID
// @Description Get a spy cat by ID
// @Tags cats
// @Accept  json
// @Produce  json
// @Param id path int true "Cat ID"
// @Success 200 {object} models.Cat
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cats/get/{id} [get]
func (app *application) getCat(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	v := validator.New()
	v.Check(id != 0, "id", validator.ErrZeroID.Error())
	if !v.Valid() {
		writeJSONValidationErrors(c, v.Errors)
		return
	}

	cat, err := app.cats.Get(c, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Cat with ID %d doesn't exist", id)})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, cat)

	slog.Info("Cat returned", "id", id)
}
