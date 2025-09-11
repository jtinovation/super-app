package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
)

type LabUseCase interface {
	FindAll(params dto.QueryParams, majorId string) (*[]domain.Lab, int64, error)
	FindByID(id string) (*domain.Lab, error)
	FindAllAsOptions(majorId string) (*[]domain.Lab, error)
	Create(dto *dto.StoreLabDTO) (*domain.Lab, error)
	Update(id string, dto *dto.UpdateLabDTO) (*domain.Lab, error)
	Delete(id string) error
}

type labUseCase struct {
	repo domain.LabRepository
}

func NewLabUseCase(repo domain.LabRepository) LabUseCase {
	return &labUseCase{repo: repo}
}

func (u *labUseCase) FindAll(params dto.QueryParams, majorId string) (*[]domain.Lab, int64, error) {
	return u.repo.FindAll(params, majorId)
}

func (u *labUseCase) FindByID(id string) (*domain.Lab, error) {
	return u.repo.FindByID(id)
}
func (u *labUseCase) FindAllAsOptions(majorId string) (*[]domain.Lab, error) {
	return u.repo.FindAllAsOptions(majorId)
}

func (u *labUseCase) Create(dto *dto.StoreLabDTO) (*domain.Lab, error) {
	lab := &domain.Lab{
		ID:          uuid.NewString(),
		Code:        dto.Code,
		Name:        dto.Name,
		MajorID:     dto.MajorID,
		EmployeeLab: make([]domain.EmployeeLab, len(dto.EmployeeLab)),
	}
	for i, emp := range dto.EmployeeLab {
		lab.EmployeeLab[i] = domain.EmployeeLab{
			ID:         uuid.NewString(),
			LabID:      lab.ID,
			EmployeeID: emp.EmployeeID,
			IsHeadLab:  emp.IsHeadLab,
			Period:     emp.Period,
		}
	}
	return u.repo.Create(lab)
}

func (u *labUseCase) Update(id string, dto *dto.UpdateLabDTO) (*domain.Lab, error) {
	lab := &domain.Lab{
		Code:    dto.Code,
		Name:    dto.Name,
		MajorID: dto.MajorID,
	}
	return u.repo.Update(id, lab)
}

func (u *labUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}
