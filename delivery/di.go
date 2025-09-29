package delivery

import (
	"jti-super-app-go/config"
	"jti-super-app-go/internal/handler"
	"jti-super-app-go/internal/repository"
	"jti-super-app-go/internal/service"
	"jti-super-app-go/internal/usecase"

	"gorm.io/gorm"
)

type Container struct {
	AuthHandler         *handler.AuthHandler
	EmployeeHandler     *handler.EmployeeHandler
	LabHandler          *handler.LabHandler
	MajorHandler        *handler.MajorHandler
	SemesterHandler     *handler.SemesterHandler
	SessionHandler      *handler.SessionHandler
	StudentHandler      *handler.StudentHandler
	StudyProgramHandler *handler.StudyProgramHandler
	SubjectHandler      *handler.SubjectHandler
	GoogleAuthService   service.GoogleAuthService
	OauthClientHandler  *handler.OauthClientHandler
	OauthHandler        *handler.OauthHandler
	PermissionHandler   *handler.PermissionHandler
	RoleHandler         *handler.RoleHandler
	UserHandler         *handler.UserHandler
}

func InitContainer(db *gorm.DB, jwtService service.JWTService) *Container {
	emailService := service.NewEmailService(config.AppConfig.Email)
	employeeRepo := repository.NewEmployeeRepository(db)
	googleAuthService := service.NewGoogleAuthService(config.AppConfig)
	studentRepo := repository.NewStudentRepository(db)

	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	authUC := usecase.NewAuthUseCase(authRepo, userRepo, employeeRepo, studentRepo, passwordResetRepo, jwtService, emailService)
	authHandler := handler.NewAuthHandler(authUC, googleAuthService)

	employeeUC := usecase.NewEmployeeUseCase(db, employeeRepo, userRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeUC)

	labRepo := repository.NewLabRepository(db)
	labUC := usecase.NewLabUseCase(labRepo)
	labHandler := handler.NewLabHandler(labUC)

	majorRepo := repository.NewMajorRepository(db)
	majorUC := usecase.NewMajorUseCase(majorRepo)
	majorHandler := handler.NewMajorHandler(majorUC)

	permissionRepo := repository.NewPermissionRepository(db)
	permissionUC := usecase.NewPermissionUseCase(permissionRepo)
	permissionHandler := handler.NewPermissionHandler(permissionUC)

	roleRepo := repository.NewRoleRepository(db)
	roleUC := usecase.NewRoleUseCase(roleRepo)
	roleHandler := handler.NewRoleHandler(roleUC)

	semesterRepo := repository.NewSemesterRepository(db)
	semesterUC := usecase.NewSemesterUseCase(semesterRepo)
	semesterHandler := handler.NewSemesterHandler(semesterUC)

	sessionRepo := repository.NewSessionRepository(db)
	sessionUC := usecase.NewSessionUseCase(sessionRepo)
	sessionHandler := handler.NewSessionHandler(sessionUC)

	studentSemesterRepo := repository.NewStudentSemesterRepository(db)
	studentUC := usecase.NewStudentUseCase(db, studentRepo, userRepo, studentSemesterRepo)
	studentHandler := handler.NewStudentHandler(studentUC)

	studyProgramRepo := repository.NewStudyProgramRepository(db)
	studyProgramUC := usecase.NewStudyProgramUseCase(studyProgramRepo)
	studyProgramHandler := handler.NewStudyProgramHandler(studyProgramUC)

	subjectSemesterRepo := repository.NewSubjectSemesterRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	subjectUC := usecase.NewSubjectUseCase(subjectRepo, subjectSemesterRepo)
	subjectHandler := handler.NewSubjectHandler(subjectUC)

	oauthClientRepo := repository.NewOauthClientRepository(db)
	oauthClientUC := usecase.NewOauthClientUseCase(oauthClientRepo)
	oauthClientHandler := handler.NewOauthClientHandler(oauthClientUC)

	oauthUsecase := usecase.NewOauthUsecase()
	oauthHandler := handler.NewOauthHandler(oauthClientUC, oauthUsecase, authUC)

	userUC := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUC)

	return &Container{
		AuthHandler:         authHandler,
		EmployeeHandler:     employeeHandler,
		LabHandler:          labHandler,
		MajorHandler:        majorHandler,
		SemesterHandler:     semesterHandler,
		SessionHandler:      sessionHandler,
		StudentHandler:      studentHandler,
		StudyProgramHandler: studyProgramHandler,
		SubjectHandler:      subjectHandler,
		OauthClientHandler:  oauthClientHandler,
		OauthHandler:        oauthHandler,
		PermissionHandler:   permissionHandler,
		RoleHandler:         roleHandler,
		UserHandler:         userHandler,
	}
}
