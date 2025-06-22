package delivery

import (
	"jti-super-app-go/delivery/middleware"
	"jti-super-app-go/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, c *Container, jwtService service.JWTService) {
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", c.AuthHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(jwtService), c.AuthHandler.Logout)
		}

		classes := api.Group("/classes").Use(middleware.AuthMiddleware(jwtService))
		{
			classes.GET("", c.ClassHandler.FindAll)
			classes.GET("/:id", c.ClassHandler.FindByID)
			classes.GET("/options", c.ClassHandler.FindAllAsOptions)
			classes.POST("", c.ClassHandler.Create)
			classes.PUT("/:id", c.ClassHandler.Update)
			classes.DELETE("/:id", c.ClassHandler.Delete)
		}

		employees := api.Group("/employees").Use(middleware.AuthMiddleware(jwtService))
		{
			employees.GET("", c.EmployeeHandler.FindAll)
			employees.GET("/options", c.EmployeeHandler.FindAllAsOptions)
			employees.GET("/:id", c.EmployeeHandler.FindByID)
			employees.POST("", c.EmployeeHandler.Create)
			employees.POST("/:id/update", c.EmployeeHandler.Update)
			employees.DELETE("/:id", c.EmployeeHandler.Delete)
		}

		majors := api.Group("/majors").Use(middleware.AuthMiddleware(jwtService))
		{
			majors.GET("", c.MajorHandler.FindAll)
			majors.GET("/:id", c.MajorHandler.FindByID)
			majors.GET("/options", c.MajorHandler.FindAllAsOptions)
			majors.POST("", c.MajorHandler.Create)
			majors.PUT("/:id", c.MajorHandler.Update)
			majors.DELETE("/:id", c.MajorHandler.Delete)
		}

		studyPrograms := api.Group("/study-programs").Use(middleware.AuthMiddleware(jwtService))
		{
			studyPrograms.GET("", c.StudyProgramHandler.FindAll)
			studyPrograms.GET("/:id", c.StudyProgramHandler.FindByID)
			studyPrograms.GET("/options", c.StudyProgramHandler.FindAllAsOptions)
			studyPrograms.POST("", c.StudyProgramHandler.Create)
			studyPrograms.PUT("/:id", c.StudyProgramHandler.Update)
			studyPrograms.DELETE("/:id", c.StudyProgramHandler.Delete)
		}
	}
}
