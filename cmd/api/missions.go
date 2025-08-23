package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"spy-cat-agency/internal/missions"
	"spy-cat-agency/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new mission
// @Description Create a new mission
// @Tags missions
// @Accept  json
// @Produce  json
// @Param mission body models.Mission true "Mission object"
// @Success 201 {object} models.Mission
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/create [post]
func (app *application) createMission(c *gin.Context) {
	var mission models.Mission

	if err := c.ShouldBindJSON(&mission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if validationErrors, ok := app.missions.ValidateMission(&mission); !ok {
		returnValidationErrors(c, validationErrors)
		return
	}

	newMission, err := app.missions.Create(c, &mission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, newMission)

	slog.Info("Mission created", "id", newMission.ID)
}

// @Summary Delete a mission
// @Description Delete a mission by ID
// @Tags missions
// @Accept  json
// @Produce  json
// @Param id path int true "Mission ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/delete/{id} [delete]
func (app *application) deleteMission(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	if err := app.missions.Delete(c, id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Mission with ID %d doesn't exist", id)})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"info": "success"})

	slog.Info("Mission deleted", "id", id)
}

// @Summary Complete a mission
// @Description Mark a mission as completed
// @Tags missions
// @Accept  json
// @Produce  json
// @Param id path int true "Mission ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/complete/{id} [put]
func (app *application) completeMission(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	if err := app.missions.UpdateAsCompleted(c, id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Mission with ID %d doesn't exist", id)})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"info": "success"})

	slog.Info("Mission marked as completed", "id", id)
}

// @Summary Update target notes
// @Description Update the notes for a target
// @Tags missions
// @Accept  json
// @Produce  json
// @Param target body models.Target true "Target object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/update_notes [put]
func (app *application) updateTargetNotes(c *gin.Context) {
	var target models.Target

	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target body"})
		return
	}

	if target.ID == 0 || target.Notes == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid target ID or notes"})
		return
	}

	if err := app.missions.UpdateTargetNotes(c, target.ID, target.Notes); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Target with ID %d doesn't exist", target.ID)})
			return
		case errors.Is(err, missions.ErrMIssionCompleted):
			c.JSON(http.StatusNotFound, gin.H{"error": "cannot update, mission is already completed"})
			return
		case errors.Is(err, missions.ErrTargetCompleted):
			c.JSON(http.StatusNotFound, gin.H{"error": "cannot update, target is already completed"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"info": "success"})

	slog.Info("Target notes updated", "id", target.ID)
}

// @Summary Delete a target
// @Description Delete a target by ID
// @Tags missions
// @Accept  json
// @Produce  json
// @Param id path int true "Target ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/delete_target/{id} [delete]
func (app *application) deleteTarget(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	if err := app.missions.DeleteTarget(c, id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Target with ID %d doesn't exist", id)})
			return
		case errors.Is(err, missions.ErrTargetCompleted):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update, target is already completed"})
			return
		case errors.Is(err, missions.ErrMissionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Associated mission not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"info": "success"})

	slog.Info("Target deleted", "id", id)
}

// @Summary Add targets to a mission
// @Description Add new targets to an existing mission
// @Tags missions
// @Accept  json
// @Produce  json
// @Param mission body models.Mission true "Mission object with new targets"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/add_targets [put]
func (app *application) addTargets(c *gin.Context) {
	var mission models.Mission

	if err := c.ShouldBindJSON(&mission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission body"})
		return
	}

	if validationErrors, ok := app.missions.ValidateMission(&mission); !ok {
		returnValidationErrors(c, validationErrors)
		return
	}

	if len(mission.Targets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no targets provided"})
	}

	newTargets, err := app.missions.AddTargets(c, mission.ID, mission.Targets)
	if err != nil {
		switch {
		case errors.Is(err, missions.ErrMIssionCompleted):
			c.JSON(http.StatusNotFound, gin.H{"error": "Unable to add targets, mission is already completed"})
			return
		case errors.Is(err, missions.ErrTooManyTargets):
			c.JSON(http.StatusNotFound, gin.H{"error": "Unable to add targets, too many targets"})
			return
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Mission with ID %d doesn't exist", mission.ID)})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"info": "success", "added targets": newTargets})

	slog.Info("Added new targets to mission", "id", mission.ID, "new targets", len(newTargets))
}

// @Summary Assign a cat to a mission
// @Description Assign a spy cat to an existing mission
// @Tags missions
// @Accept  json
// @Produce  json
// @Param mission body models.Mission true "Mission object with cat ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/assign [put]
func (app *application) assignCat(c *gin.Context) {
	var mission models.Mission

	if err := c.ShouldBindJSON(&mission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission body"})
		return
	}

	if validationErrors, ok := app.missions.ValidateMission(&mission); !ok {
		returnValidationErrors(c, validationErrors)
		return
	}

	if err := app.missions.AssignCat(c, mission.ID, *mission.CatID); err != nil {
		switch {
		case errors.Is(err, missions.ErrMissionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Mission with id %d doesn't exist", mission.ID)})
			return
		case errors.Is(err, missions.ErrCatNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cat with id %d doesn't exist", *mission.CatID)})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"info": "success"})

	slog.Info("Cat assigned to a mission", "cat id", mission.CatID, "mission id", mission.ID)
}

// @Summary List all missions
// @Description Get a list of all missions
// @Tags missions
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Mission
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/list [get]
func (app *application) listMissions(c *gin.Context) {
	missions, err := app.missions.List(c)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": "No missions found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, missions)

	slog.Info("Missions returned", "missions", len(*missions))
}

// @Summary Get a mission by ID
// @Description Get a mission by ID
// @Tags missions
// @Accept  json
// @Produce  json
// @Param id path int true "Mission ID"
// @Success 200 {object} models.Mission
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/get/{id} [get]
func (app *application) getMission(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}
	mission, err := app.missions.Get(c, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Mission with ID %d doesn't exist", id)})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, mission)

	slog.Info("Mission returned", "mission ID", id)
}
