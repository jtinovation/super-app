package repository

import (
	"gorm.io/gorm"
)

type employeeLabRepository struct {
	db *gorm.DB
}

func NewEmployeeLabRepository(db *gorm.DB) *employeeLabRepository {
	return &employeeLabRepository{db: db}
}

// func (r *employeeLabRepository) FindAll(params dto.QueryParams, majorId string) (*[]domain.EmployeeLab, int64, error) {
// 	var employeeLabs []domain.EmployeeLab
// 	var totalRows int64

// 	query := r.db.Model(&domain.EmployeeLab{}).Where("major_id = ?", majorId)

// 	if params.Search != "" {
// 		query = query.Where("name ILIKE ?", "%"+params.Search+"%")
// 	}

// 	if err := query.Count(&totalRows).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	if err := query.Offset(params.Offset).Limit(params.Limit).Find(&employeeLabs).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	return &employeeLabs, totalRows, nil
// }
