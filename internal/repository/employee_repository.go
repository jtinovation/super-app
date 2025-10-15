package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) domain.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) FindAll(params dto.QueryParams, position string, majorId string) (*[]domain.Employee, int64, error) {
	var employees []domain.Employee
	var totalRows int64

	query := r.db.Model(&domain.Employee{}).
		Select("m_employee.*, m_user.name as name, m_user.email as email, m_user.img_path as img_path, m_user.img_name as img_name").
		Joins("LEFT JOIN m_user ON m_user.id = m_employee.m_user_id")

	if params.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(params.Search))
		query = query.Where(
			r.db.
				Where("LOWER(m_employee.nip) LIKE ?", searchQuery).
				Or("LOWER(m_user.name) LIKE ?", searchQuery).
				Or("LOWER(m_user.email) LIKE ?", searchQuery).
				Or("LOWER(m_employee.position) LIKE ?", searchQuery),
		)
	}

	if position != "" {
		query = query.Where("m_employee.position = ?", position)
	}

	if majorId != "" {
		query = query.Where("m_employee.m_major_id = ?", majorId)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := fmt.Sprintf("%s %s", params.Sort, params.Order)
		query = query.Order(sortOrder)
	} else {
		query = query.Order("m_user.name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return &employees, totalRows, nil
}

func (r *employeeRepository) FindByID(id string) (*domain.Employee, error) {
	var employee domain.Employee
	if err := r.db.Preload("User").Preload("Major").First(&employee, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) FindByUserID(userID string) (*domain.Employee, error) {
	var employee domain.Employee
	if err := r.db.Preload("Major").Preload("StudyProgram").First(&employee, "m_user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) FindAllAsOptions(position string, majorId string, studyProgramId string) (*[]domain.Employee, error) {
	var employees []domain.Employee
	query := r.db.Model(&domain.Employee{}).
		Select("m_employee.id, m_user.name").
		Joins("LEFT JOIN m_user ON m_user.id = m_employee.m_user_id")

	if position != "" {
		query = query.Where("position = ?", position)
	}
	if majorId != "" {
		query = query.Where("m_major_id = ?", majorId)
	}
	if studyProgramId != "" {
		query = query.Where("m_study_program_id = ?", studyProgramId)
	}

	if err := query.Find(&employees).Error; err != nil {
		return nil, err
	}
	return &employees, nil
}

func (r *employeeRepository) Create(employee *domain.Employee) (*domain.Employee, error) {
	if err := r.db.Create(employee).Error; err != nil {
		return nil, err
	}
	return employee, nil
}

func (r *employeeRepository) Update(id string, employee *domain.Employee) (*domain.Employee, error) {
	var existingEmployee domain.Employee
	if err := r.db.First(&existingEmployee, "id = ?", id).Error; err != nil {
		return nil, err // record not found
	}

	if err := r.db.Model(&existingEmployee).Updates(employee).Error; err != nil {
		return nil, err
	}
	return &existingEmployee, nil
}

func (r *employeeRepository) Delete(id string) error {
	if err := r.db.First(&domain.Employee{}, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&domain.Employee{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
