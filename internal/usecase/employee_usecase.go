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

type EmployeeUseCase interface {
	FindAll(params dto.QueryParams, position string, majorId string) (*[]domain.Employee, int64, error)
	FindByID(id string) (*domain.Employee, error)
	FindAllAsOptions(position string, majorId string, studyProgramId string) (*[]domain.Employee, error)
	Create(dto *dto.StoreEmployeeDTO) (*domain.Employee, error)
	Update(id string, dto *dto.UpdateEmployeeDTO) (*domain.Employee, error)
	Delete(id string) error
}

type employeeUseCase struct {
	db       *gorm.DB
	empRepo  domain.EmployeeRepository
	userRepo domain.UserRepository
}

func NewEmployeeUseCase(db *gorm.DB, empRepo domain.EmployeeRepository, userRepo domain.UserRepository) EmployeeUseCase {
	return &employeeUseCase{
		db:       db,
		empRepo:  empRepo,
		userRepo: userRepo,
	}
}

func (u *employeeUseCase) FindAll(params dto.QueryParams, position string, majorId string) (*[]domain.Employee, int64, error) {
	return u.empRepo.FindAll(params, position, majorId)
}

func (u *employeeUseCase) FindByID(id string) (*domain.Employee, error) {
	return u.empRepo.FindByID(id)
}

func (u *employeeUseCase) FindAllAsOptions(position string, majorId string, studyProgramId string) (*[]domain.Employee, error) {
	return u.empRepo.FindAllAsOptions(position, majorId, studyProgramId)
}

func (u *employeeUseCase) Create(payload *dto.StoreEmployeeDTO) (*domain.Employee, error) {
	var newEmployee *domain.Employee
	imgPath := constants.EMPLOYEE_PATH
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NIP), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var birthDate *time.Time
	if payload.BirthDate != nil {
		parsedTime, err := time.Parse("2006-01-02", *payload.BirthDate)
		if err == nil {
			birthDate = &parsedTime
		}
	}

	newUser := &domain.User{
		ID:          uuid.NewString(),
		Name:        payload.Name,
		Email:       payload.Email,
		Password:    string(hashedPassword),
		Status:      constants.StatusActive,
		Gender:      payload.Gender,
		Religion:    payload.Religion,
		BirthDate:   birthDate,
		BirthPlace:  payload.BirthPlace,
		Address:     payload.Address,
		PhoneNumber: payload.PhoneNumber,
		Nationality: payload.Nationality,
		ImgPath:     &imgPath,
		ImgName:     &imgName,
	}

	createdUser, err := u.userRepo.Create(newUser)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	employee := &domain.Employee{
		ID:       uuid.NewString(),
		UserID:   createdUser.ID,
		MajorID:  payload.MajorID,
		Nip:      payload.NIP,
		Position: payload.Position,
	}

	newEmployee, err = u.empRepo.Create(employee)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return newEmployee, nil
}

func (u *employeeUseCase) Update(id string, payload *dto.UpdateEmployeeDTO) (*domain.Employee, error) {
	imgPath := constants.EMPLOYEE_PATH

	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	employeeToUpdate, err := u.empRepo.FindByID(id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var birthDate *time.Time
	if payload.BirthDate != nil {
		parsedTime, err := time.Parse("2006-01-02", *payload.BirthDate)
		if err == nil {
			birthDate = &parsedTime
		}
	}

	userUpdateData := &domain.User{
		Name:        payload.Name,
		Email:       payload.Email,
		Gender:      payload.Gender,
		Religion:    payload.Religion,
		BirthDate:   birthDate,
		BirthPlace:  payload.BirthPlace,
		Address:     payload.Address,
		PhoneNumber: payload.PhoneNumber,
		Nationality: payload.Nationality,
	}

	if payload.Avatar != nil {
		extension := filepath.Ext(payload.Avatar.Filename)
		imgName := fmt.Sprintf("%s%s", uuid.NewString(), extension)

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

		if employeeToUpdate.User.ImgPath != nil && employeeToUpdate.User.ImgName != nil && *employeeToUpdate.User.ImgName != constants.DEFAULT_AVATAR {
			oldObject := fmt.Sprintf("%s/%s", *employeeToUpdate.User.ImgPath, *employeeToUpdate.User.ImgName)
			helper.DeleteFile(config.AppConfig.Minio.Bucket, oldObject)
		}

		userUpdateData.ImgPath = &imgPath
		userUpdateData.ImgName = &imgName
	}

	_, err = u.userRepo.Update(employeeToUpdate.UserID, userUpdateData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	employeeUpdateData := &domain.Employee{
		MajorID:  payload.MajorID,
		Nip:      payload.NIP,
		Position: payload.Position,
	}

	updatedEmployee, err := u.empRepo.Update(id, employeeUpdateData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return updatedEmployee, nil
}

func (u *employeeUseCase) Delete(id string) error {
	return u.empRepo.Delete(id)
}
