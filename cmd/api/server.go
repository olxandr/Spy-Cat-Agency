package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() *gin.Engine {
	r := gin.Default()
	r.Use(loggingMiddleware())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Cats
	cats := r.Group("/cats")
	{
		cats.POST("/create", app.createCat)
		cats.DELETE("/remove/:id", app.removeCat)
		cats.PUT("/update_salary", app.updateCatsSalary)
		cats.GET("/list", app.listCats)
		cats.GET("/get/:id", app.getCat)
	}

	// Missions
	missions := r.Group("/missions")
	{
		missions.POST("/create", app.createMission)
		missions.DELETE("/delete/:id", app.deleteMission)
		missions.PUT("/complete/:id", app.completeMission)
		missions.PUT("/update_notes", app.updateTargetNotes)
		missions.DELETE("/delete_target/:id", app.deleteTarget)
		missions.PUT("/add_targets", app.addTargets)
		missions.PUT("/assign", app.assignCat)
		missions.GET("/list", app.listMissions)
		missions.GET("/get/:id", app.getMission)
	}

	r.GET("/healthcheck", app.healthcheck)

	return r
}

func (app *application) healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
