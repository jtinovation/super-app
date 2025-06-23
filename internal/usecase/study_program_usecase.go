package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
)

type StudyProgramUseCase interface {
	FindAll(params dto.QueryParams, majorId string) (*[]domain.StudyProgram, int64, error)
	FindByID(id string) (*domain.StudyProgram, error)
	FindAllAsOptions(majorId string) (*[]domain.StudyProgram, error)
	Create(dto *dto.StoreStudyProgramDTO) (*domain.StudyProgram, error)
	Update(id string, dto *dto.UpdateStudyProgramDTO) (*domain.StudyProgram, error)
	Delete(id string) error
}

type studyProgramUseCase struct {
	repo domain.StudyProgramRepository
}

func NewStudyProgramUseCase(repo domain.StudyProgramRepository) StudyProgramUseCase {
	return &studyProgramUseCase{repo: repo}
}

func (u *studyProgramUseCase) FindAll(params dto.QueryParams, majorId string) (*[]domain.StudyProgram, int64, error) {
	return u.repo.FindAll(params, majorId)
}

func (u *studyProgramUseCase) FindByID(id string) (*domain.StudyProgram, error) {
	return u.repo.FindByID(id)
}
func (u *studyProgramUseCase) FindAllAsOptions(majorId string) (*[]domain.StudyProgram, error) {
	return u.repo.FindAllAsOptions(majorId)
}

func (u *studyProgramUseCase) Create(dto *dto.StoreStudyProgramDTO) (*domain.StudyProgram, error) {
	studyProgram := &domain.StudyProgram{
		ID:      uuid.NewString(),
		Code:    dto.Code,
		Name:    dto.Name,
		MajorID: dto.MajorID,
	}
	return u.repo.Create(studyProgram)
}

func (u *studyProgramUseCase) Update(id string, dto *dto.UpdateStudyProgramDTO) (*domain.StudyProgram, error) {
	studyProgram := &domain.StudyProgram{
		Code:    dto.Code,
		Name:    dto.Name,
		MajorID: dto.MajorID,
	}
	return u.repo.Update(id, studyProgram)
}

func (u *studyProgramUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}
