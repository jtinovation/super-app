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
	MajorHandler        *handler.MajorHandler
	SemesterHandler     *handler.SemesterHandler
	SessionHandler      *handler.SessionHandler
	StudentHandler      *handler.StudentHandler
	StudyProgramHandler *handler.StudyProgramHandler
	SubjectHandler      *handler.SubjectHandler
	GoogleAuthService   service.GoogleAuthService
}

func InitContainer(db *gorm.DB, jwtService service.JWTService) *Container {
	emailService := service.NewEmailService(config.AppConfig.Email)
	googleAuthService := service.NewGoogleAuthService(config.AppConfig)

	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	authUC := usecase.NewAuthUseCase(authRepo, userRepo, passwordResetRepo, jwtService, emailService)
	authHandler := handler.NewAuthHandler(authUC, googleAuthService)

	employeeRepo := repository.NewEmployeeRepository(db)
	employeeUC := usecase.NewEmployeeUseCase(db, employeeRepo, userRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeUC)

	majorRepo := repository.NewMajorRepository(db)
	majorUC := usecase.NewMajorUseCase(majorRepo)
	majorHandler := handler.NewMajorHandler(majorUC)

	semesterRepo := repository.NewSemesterRepository(db)
	semesterUC := usecase.NewSemesterUseCase(semesterRepo)
	semesterHandler := handler.NewSemesterHandler(semesterUC)

	sessionRepo := repository.NewSessionRepository(db)
	sessionUC := usecase.NewSessionUseCase(sessionRepo)
	sessionHandler := handler.NewSessionHandler(sessionUC)

	studentRepo := repository.NewStudentRepository(db)
	studentUC := usecase.NewStudentUseCase(db, studentRepo, userRepo)
	studentHandler := handler.NewStudentHandler(studentUC)

	studyProgramRepo := repository.NewStudyProgramRepository(db)
	studyProgramUC := usecase.NewStudyProgramUseCase(studyProgramRepo)
	studyProgramHandler := handler.NewStudyProgramHandler(studyProgramUC)

	subjectSemesterRepo := repository.NewSubjectSemesterRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	subjectUC := usecase.NewSubjectUseCase(subjectRepo, subjectSemesterRepo)
	subjectHandler := handler.NewSubjectHandler(subjectUC)

	return &Container{
		AuthHandler:         authHandler,
		EmployeeHandler:     employeeHandler,
		MajorHandler:        majorHandler,
		SemesterHandler:     semesterHandler,
		SessionHandler:      sessionHandler,
		StudentHandler:      studentHandler,
		StudyProgramHandler: studyProgramHandler,
		SubjectHandler:      subjectHandler,
	}
}
