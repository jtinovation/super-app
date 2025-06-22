package delivery

import (
	"jti-super-app-go/internal/handler"
	"jti-super-app-go/internal/repository"
	"jti-super-app-go/internal/service"
	"jti-super-app-go/internal/usecase"

	"gorm.io/gorm"
)

type Container struct {
	AuthHandler         *handler.AuthHandler
	ClassHandler        *handler.ClassHandler
	EmployeeHandler     *handler.EmployeeHandler
	MajorHandler        *handler.MajorHandler
	StudyProgramHandler *handler.StudyProgramHandler
}

func InitContainer(db *gorm.DB, jwtService service.JWTService) *Container {
	authRepo := repository.NewAuthRepository(db)
	authUC := usecase.NewAuthUseCase(authRepo, jwtService)
	authHandler := handler.NewAuthHandler(authUC)

	classRepo := repository.NewClassRepository(db)
	classUC := usecase.NewClassUseCase(classRepo)
	classHandler := handler.NewClassHandler(classUC)

	userRepo := repository.NewUserRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	employeeUC := usecase.NewEmployeeUseCase(db, employeeRepo, userRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeUC)

	majorRepo := repository.NewMajorRepository(db)
	majorUC := usecase.NewMajorUseCase(majorRepo)
	majorHandler := handler.NewMajorHandler(majorUC)

	studyProgramRepo := repository.NewStudyProgramRepository(db)
	studyProgramUC := usecase.NewStudyProgramUseCase(studyProgramRepo)
	studyProgramHandler := handler.NewStudyProgramHandler(studyProgramUC)

	return &Container{
		AuthHandler:         authHandler,
		ClassHandler:        classHandler,
		EmployeeHandler:     employeeHandler,
		MajorHandler:        majorHandler,
		StudyProgramHandler: studyProgramHandler,
	}
}
