package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
)

type ClassUseCase interface {
	FindAll(params dto.QueryParams, studyProgramId string, majorId string) (*[]domain.Class, int64, error)
	FindByID(id string) (*domain.Class, error)
	FindAllAsOptions(studyProgramId string) (*[]domain.Class, error)
	Create(dto *dto.StoreClassDTO) (*domain.Class, error)
	Update(id string, dto *dto.UpdateClassDTO) (*domain.Class, error)
	Delete(id string) error
}

type classUseCase struct {
	repo domain.ClassRepository
}

func NewClassUseCase(repo domain.ClassRepository) ClassUseCase {
	return &classUseCase{repo: repo}
}

func (u *classUseCase) FindAll(params dto.QueryParams, studyProgramId string, majorId string) (*[]domain.Class, int64, error) {
	return u.repo.FindAll(params, studyProgramId, majorId)
}

func (u *classUseCase) FindByID(id string) (*domain.Class, error) {
	return u.repo.FindByID(id)
}
func (u *classUseCase) FindAllAsOptions(studyProgramId string) (*[]domain.Class, error) {
	return u.repo.FindAllAsOptions(studyProgramId)
}

func (u *classUseCase) Create(dto *dto.StoreClassDTO) (*domain.Class, error) {
	class := &domain.Class{
		ID:             uuid.NewString(),
		Code:           dto.Code,
		Name:           dto.Name,
		StudyProgramID: dto.StudyProgramID,
	}
	return u.repo.Create(class)
}

func (u *classUseCase) Update(id string, dto *dto.UpdateClassDTO) (*domain.Class, error) {
	class := &domain.Class{
		Code:           dto.Code,
		Name:           dto.Name,
		StudyProgramID: dto.StudyProgramID,
	}
	return u.repo.Update(id, class)
}

func (u *classUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}
