package usecase

import (
	"fmt"
	"jti-super-app-go/config"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/pkg/constants"
	"jti-super-app-go/pkg/helper"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StudentUseCase interface {
	FindAll(params dto.QueryParams) (*[]domain.Student, int64, error)
	Create(payload *dto.StoreStudentDTO) (*domain.Student, error)
}

type studentUseCase struct {
	db                  *gorm.DB
	studentRepo         domain.StudentRepository
	userRepo            domain.UserRepository
	studentSemesterRepo domain.StudentSemesterRepository
}

func NewStudentUseCase(db *gorm.DB, studentRepo domain.StudentRepository, userRepo domain.UserRepository, studentSemesterRepo domain.StudentSemesterRepository) StudentUseCase {
	return &studentUseCase{
		db:                  db,
		studentRepo:         studentRepo,
		userRepo:            userRepo,
		studentSemesterRepo: studentSemesterRepo,
	}
}

func (u *studentUseCase) FindAll(params dto.QueryParams) (*[]domain.Student, int64, error) {
	return u.studentRepo.FindAll(params)
}

func (u *studentUseCase) Create(payload *dto.StoreStudentDTO) (*domain.Student, error) {
	var newStudent *domain.Student
	imgPath := constants.STUDENT_PATH
	imgName := constants.DEFAULT_AVATAR

	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if payload.Avatar != nil {
		extension := filepath.Ext(payload.Avatar.Filename)
		imgName = fmt.Sprintf("%s%s", uuid.NewString(), extension)

		file, err := payload.Avatar.Open()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		defer file.Close()

		err = helper.UploadFile(config.AppConfig.Minio.Bucket, fmt.Sprintf("%s/%s", imgPath, imgName), file, payload.Avatar.Size, payload.Avatar.Header.Get("Content-Type"))
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NIM), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var birthDate *time.Time
	if payload.BirthDate != nil {
		parsedTime, _ := time.Parse("2006-01-02", *payload.BirthDate)
		birthDate = &parsedTime
	}

	newUser := &domain.User{
		ID:          uuid.NewString(),
		Name:        payload.Name,
		Email:       payload.Email,
		Password:    string(hashedPassword),
		ImgPath:     &imgPath,
		ImgName:     &imgName,
		Status:      "ACTIVE",
		Gender:      payload.Gender,
		Religion:    payload.Religion,
		BirthPlace:  payload.BirthPlace,
		BirthDate:   birthDate,
		PhoneNumber: payload.PhoneNumber,
		Nationality: payload.Nationality,
		Address:     payload.Address,
	}
	createdUser, err := u.userRepo.Create(newUser)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create Student
	student := &domain.Student{
		ID:               uuid.NewString(),
		UserID:           createdUser.ID,
		StudentProgramID: payload.StudyProgramID,
		NIM:              payload.NIM,
		Generation:       payload.Generation,
		TuitionFee:       payload.TuitionFee,
		TuitionMethod:    payload.TuitionMethod,
		StudyProgramID:   payload.StudyProgramID,
	}

	newStudent, err = u.studentRepo.Create(student)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = u.studentSemesterRepo.StoreStudentSemester(&domain.StudentSemester{
		ID:         uuid.NewString(),
		StudentID:  newStudent.ID,
		SemesterID: payload.SemesterId,
		Class:      payload.Class,
	})

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return newStudent, nil
}
