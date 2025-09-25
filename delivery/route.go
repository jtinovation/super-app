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
			auth.GET("/email/verify/:token", c.AuthHandler.VerifyEmail)
			auth.POST("/email/resend", middleware.AuthMiddleware(jwtService), c.AuthHandler.ResendVerificationEmail)
			auth.GET("/google/login", c.AuthHandler.GoogleLogin)
			auth.GET("/google/callback", c.AuthHandler.GoogleCallback)
			auth.POST("/login", middleware.RateLimiter(), c.AuthHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(jwtService), c.AuthHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(jwtService), c.AuthHandler.Me)
			auth.POST("/password/forgot", c.AuthHandler.ForgotPassword)
			auth.POST("/password/reset", c.AuthHandler.ResetPassword)
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

		labs := api.Group("/labs").Use(middleware.AuthMiddleware(jwtService))
		{
			labs.GET("", c.LabHandler.FindAll)
			labs.GET("/:id", c.LabHandler.FindByID)
			labs.GET("/options", c.LabHandler.FindAllAsOptions)
			labs.POST("", c.LabHandler.Create)
			labs.PUT("/:id", c.LabHandler.Update)
			labs.DELETE("/:id", c.LabHandler.Delete)
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

		permissions := api.Group("/permissions").Use(middleware.AuthMiddleware(jwtService))
		{
			permissions.GET("", c.PermissionHandler.FindAll)
			permissions.GET("/options", c.PermissionHandler.FindAllAsOptions)
			permissions.GET("/:id", c.PermissionHandler.FindByID)
			permissions.POST("", c.PermissionHandler.Create)
			permissions.PUT("/:id", c.PermissionHandler.Update)
			permissions.DELETE("/:id", c.PermissionHandler.Delete)
		}

		roles := api.Group("/roles").Use(middleware.AuthMiddleware(jwtService))
		{
			roles.GET("", c.RoleHandler.FindAll)
			roles.GET("/:id", c.RoleHandler.FindByID)
			roles.POST("", c.RoleHandler.Create)
			roles.PUT("/:id", c.RoleHandler.Update)
			roles.DELETE("/:id", c.RoleHandler.Delete)
		}

		semesters := api.Group("/semesters").Use(middleware.AuthMiddleware(jwtService))
		{
			semesters.GET("", c.SemesterHandler.FindAll)
			semesters.GET("/options", c.SemesterHandler.FindAllAsOptions)
			semesters.POST("", c.SemesterHandler.Create)
			semesters.PUT("/:id", c.SemesterHandler.Update)
			semesters.DELETE("/:id", c.SemesterHandler.Delete)
			semesters.POST("/:id/setting-subjects", c.SemesterHandler.SettingSubjectSemester)
		}

		sessions := api.Group("/sessions").Use(middleware.AuthMiddleware(jwtService))
		{
			sessions.GET("", c.SessionHandler.FindAll)
			sessions.GET("/options", c.SessionHandler.FindAllAsOptions)
			sessions.POST("", c.SessionHandler.Create)
			sessions.PUT("/:id", c.SessionHandler.Update)
			sessions.DELETE("/:id", c.SessionHandler.Delete)
		}

		students := api.Group("/students").Use(middleware.AuthMiddleware(jwtService))
		{
			students.GET("", c.StudentHandler.FindAll)
			students.GET("/:id", c.StudentHandler.FindByID)
			students.POST("", c.StudentHandler.Create)
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

		subjects := api.Group("/subjects").Use(middleware.AuthMiddleware(jwtService))
		{
			subjects.GET("", c.SubjectHandler.FindAll)
			subjects.GET("/options", c.SubjectHandler.FindAllAsOptions)
			subjects.POST("", c.SubjectHandler.Create)
			subjects.PUT("/:id", c.SubjectHandler.Update)
			subjects.DELETE("/:id", c.SubjectHandler.Delete)
			subjects.GET("/lectures", c.SubjectHandler.GetLectureOnSubject)
			subjects.POST("/lectures", c.SubjectHandler.StoreLectureOnSubject)
		}

		oauthClients := api.Group("/oauth-clients")
		{
			oauthClients.GET("", c.OauthClientHandler.FindAll)
			oauthClients.GET("/:id", c.OauthClientHandler.FindByID)
			oauthClients.POST("", c.OauthClientHandler.Create)
			oauthClients.PUT("/:id", c.OauthClientHandler.Update)
			oauthClients.DELETE("/:id", c.OauthClientHandler.Delete)
		}

		oauth := api.Group("/oauth")
		{
			oauth.POST("/login", middleware.RateLimiter(), c.OauthHandler.LoginPost)
			oauth.POST("/token", c.OauthHandler.Token)
			oauth.GET("/authorize", c.OauthHandler.Authorize)
			oauth.GET("/logout", c.OauthHandler.Logout)
		}
	}

	web := router.Group("/")
	{
		web.GET("", c.OauthHandler.IndexPage)
		web.GET("/login", middleware.CSRFTokenMiddleware(), c.OauthHandler.LoginPage)
		web.GET("/auth/callback", c.OauthHandler.LoginCallback)
	}
}
