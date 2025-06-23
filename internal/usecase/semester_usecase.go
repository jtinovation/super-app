package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
)

type SemesterUseCase interface {
	FindAll(params dto.QueryParams) (*[]domain.Semester, int64, error)
	FindAllAsOptions(sessionID string) (*[]domain.Semester, error)
	Create(dto *dto.StoreSemesterDTO) (*domain.Semester, error)
	Update(id string, dto *dto.UpdateSemesterDTO) (*domain.Semester, error)
	Delete(id string) error
	SettingSubjectSemester(semesterID string, subjectIDs []string) error
}

type semesterUseCase struct {
	repo domain.SemesterRepository
}

func NewSemesterUseCase(repo domain.SemesterRepository) SemesterUseCase {
	return &semesterUseCase{repo: repo}
}

func (u *semesterUseCase) FindAll(params dto.QueryParams) (*[]domain.Semester, int64, error) {
	return u.repo.FindAll(params)
}

func (u *semesterUseCase) FindAllAsOptions(sessionID string) (*[]domain.Semester, error) {
	return u.repo.FindAllAsOptions(sessionID)
}

func (u *semesterUseCase) Create(dto *dto.StoreSemesterDTO) (*domain.Semester, error) {
	semester := &domain.Semester{
		ID:        uuid.NewString(),
		SessionID: dto.SessionID,
		Year:      dto.Year,
		Semester:  dto.Semester,
	}
	return u.repo.Create(semester)
}

func (u *semesterUseCase) Update(id string, dto *dto.UpdateSemesterDTO) (*domain.Semester, error) {
	semester := &domain.Semester{
		SessionID: dto.SessionID,
		Year:      dto.Year,
		Semester:  dto.Semester,
	}
	return u.repo.Update(id, semester)
}

func (u *semesterUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}

func (u *semesterUseCase) SettingSubjectSemester(semesterID string, subjectIDs []string) error {
	return u.repo.SettingSubjectSemester(semesterID, subjectIDs)
}
