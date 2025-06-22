// internal/usecase/student_usecase.go

package usecase

import (
	"fmt"
	"jti-super-app-go/config"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
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
	db          *gorm.DB
	studentRepo domain.StudentRepository
	userRepo    domain.UserRepository
}

func NewStudentUseCase(db *gorm.DB, studentRepo domain.StudentRepository, userRepo domain.UserRepository) StudentUseCase {
	return &studentUseCase{
		db:          db,
		studentRepo: studentRepo,
		userRepo:    userRepo,
	}
}

func (u *studentUseCase) FindAll(params dto.QueryParams) (*[]domain.Student, int64, error) {
	return u.studentRepo.FindAll(params)
}

func (u *studentUseCase) Create(payload *dto.StoreStudentDTO) (*domain.Student, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Handle avatar upload
	imgPath := "students"
	imgName := "default.png"
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

	// Create User from NIM
	email := fmt.Sprintf("%s@student.jti.polinema.ac.id", payload.NIM)
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
		ID:         uuid.NewString(),
		Name:       payload.Name,
		Email:      email,
		Password:   string(hashedPassword),
		ImgPath:    &imgPath,
		ImgName:    &imgName,
		Status:     "ACTIVE",
		Gender:     payload.Gender,
		Religion:   payload.Religion,
		BirthPlace: payload.BirthPlace,
		BirthDate:  birthDate,
	}
	createdUser, err := u.userRepo.Create(newUser)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create Student
	newStudent := &domain.Student{
		ID:            uuid.NewString(),
		UserID:        createdUser.ID,
		ClassID:       payload.ClassID,
		NIM:           payload.NIM,
		Generation:    payload.Generation,
		TuitionFee:    payload.TuitionFee,
		TuitionMethod: payload.TuitionMethod,
	}
	createdStudent, err := u.studentRepo.Create(newStudent)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return createdStudent, tx.Commit().Error
}
