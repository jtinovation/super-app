package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
)

type SubjectUseCase interface {
	FindAll(params dto.QueryParams) (*[]domain.Subject, int64, error)
	FindAllAsOptions(studyProgramID, semesterID string) (*[]domain.Subject, error)
	Create(payload *dto.StoreSubjectDTO) (*domain.Subject, error)
	Update(id string, payload *dto.UpdateSubjectDTO) (*domain.Subject, error)
	Delete(id string) error
	GetLectureOnSubject(studyProgramID, semesterID string) (*[]domain.SubjectSemester, error)
	StoreLectureOnSubject(data []dto.LectureMappingDTO) error
}

type subjectUseCase struct {
	subjectRepo         domain.SubjectRepository
	subjectSemesterRepo domain.SubjectSemesterRepository
}

func NewSubjectUseCase(subjectRepo domain.SubjectRepository, subjectSemesterRepo domain.SubjectSemesterRepository) SubjectUseCase {
	return &subjectUseCase{subjectRepo: subjectRepo, subjectSemesterRepo: subjectSemesterRepo}
}

func (u *subjectUseCase) FindAll(params dto.QueryParams) (*[]domain.Subject, int64, error) {
	return u.subjectRepo.FindAll(params)
}

func (u *subjectUseCase) FindAllAsOptions(studyProgramID, semesterID string) (*[]domain.Subject, error) {
	return u.subjectRepo.FindAllAsOptions(studyProgramID, semesterID)
}

func (u *subjectUseCase) Create(payload *dto.StoreSubjectDTO) (*domain.Subject, error) {
	status := "ACTIVE"
	if payload.Status != nil {
		status = *payload.Status
	}
	subject := &domain.Subject{
		ID:             uuid.NewString(),
		StudyProgramID: payload.StudyProgramID,
		Code:           payload.Code,
		Name:           payload.Name,
		Status:         status,
	}
	return u.subjectRepo.Create(subject)
}

func (u *subjectUseCase) Update(id string, payload *dto.UpdateSubjectDTO) (*domain.Subject, error) {
	subject := &domain.Subject{
		StudyProgramID: payload.StudyProgramID,
		Code:           payload.Code,
		Name:           payload.Name,
		Status:         *payload.Status,
	}
	return u.subjectRepo.Update(id, subject)
}

func (u *subjectUseCase) Delete(id string) error { return u.subjectRepo.Delete(id) }

func (u *subjectUseCase) GetLectureOnSubject(studyProgramID, semesterID string) (*[]domain.SubjectSemester, error) {
	return u.subjectSemesterRepo.GetLectureOnSubject(studyProgramID, semesterID)
}

func (u *subjectUseCase) StoreLectureOnSubject(data []dto.LectureMappingDTO) error {
	return u.subjectSemesterRepo.StoreLectureOnSubject(data)
}
