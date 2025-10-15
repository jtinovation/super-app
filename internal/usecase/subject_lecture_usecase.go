package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
)

type SubjectLectureUseCase interface {
	FindAll(params dto.QueryParams) (*[]domain.SubjectLecture, int64, error)
}

type subjectLectureUseCase struct {
	subjectLectureRepo domain.SubjectLectureRepository
}

func NewSubjectLectureUseCase(subjectLectureRepo domain.SubjectLectureRepository) SubjectLectureUseCase {
	return &subjectLectureUseCase{subjectLectureRepo: subjectLectureRepo}
}

func (u *subjectLectureUseCase) FindAll(params dto.QueryParams) (*[]domain.SubjectLecture, int64, error) {
	return u.subjectLectureRepo.FindAll(params)
}
