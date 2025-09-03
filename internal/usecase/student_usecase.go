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
	FindByID(id string) (*domain.Student, error)
	Create(payload *dto.StoreStudentDTO) (*domain.Student, error)
	Update(id string, payload *dto.UpdateStudentDTO) (*domain.Student, error)
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

func (u *studentUseCase) FindByID(id string) (*domain.Student, error) {
	student, err := u.studentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if student == nil {
		return nil, fmt.Errorf("student with ID %s not found", id)
	}

	return student, nil
}

func (u *studentUseCase) Update(id string, payload *dto.UpdateStudentDTO) (*domain.Student, error) {
	imgPath := constants.STUDENT_PATH
	imgName := constants.DEFAULT_AVATAR

	student, err := u.studentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if student == nil {
		return nil, fmt.Errorf("student with ID %s not found", id)
	}

	tx := u.db.Begin()

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

	var birthDate *time.Time
	if payload.BirthDate != nil {
		parsedTime, _ := time.Parse("2006-01-02", *payload.BirthDate)
		birthDate = &parsedTime
	}

	// Update User
	student.User.Name = payload.Name
	student.User.Email = payload.Email
	student.User.ImgPath = &imgPath
	student.User.ImgName = &imgName
	student.User.Status = *payload.Status
	student.User.Gender = payload.Gender
	student.User.Religion = payload.Religion
	student.User.BirthPlace = payload.BirthPlace
	student.User.BirthDate = birthDate
	student.User.PhoneNumber = payload.PhoneNumber
	student.User.Nationality = payload.Nationality
	student.User.Address = payload.Address

	updatedUser, err := u.userRepo.Update(student.User.ID, &student.User)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// update Student
	student.User = *updatedUser
	student.NIM = payload.NIM
	student.Generation = payload.Generation
	student.TuitionFee = payload.TuitionFee
	student.TuitionMethod = payload.TuitionMethod
	student.StudyProgramID = payload.StudyProgramID
	updatedStudent, err := u.studentRepo.Update(id, student)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return updatedStudent, nil
}
