package repository

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(params dto.QueryParams) (*[]domain.User, int64, error) {
	var users []domain.User
	var totalRows int64

	query := r.db.Model(&domain.User{}).Preload("Roles")

	if params.Search != "" {
		searchQuery := "%" + params.Search + "%"

		// Join with roles for searching by role name
		query = query.Joins("LEFT JOIN model_has_roles ON model_has_roles.model_uuid = m_user.id").
			Joins("LEFT JOIN roles ON roles.uuid = model_has_roles.role_id").
			Where(
				r.db.
					Where("m_user.name LIKE ? OR m_user.email LIKE ? OR m_user.status LIKE ?", searchQuery, searchQuery, searchQuery).
					Or("roles.name LIKE ?", searchQuery),
			).
			Group("m_user.id")
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := params.Sort + " " + params.Order
		query = query.Order(sortOrder)
	} else {
		query = query.Order("name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return &users, totalRows, nil
}

func (r *userRepository) Create(user *domain.User) (*domain.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(id string, user *domain.User) (*domain.User, error) {
	var existingUser domain.User
	if err := r.db.First(&existingUser, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingUser).Updates(user).Error; err != nil {
		return nil, err
	}
	return &existingUser, nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("Roles.Permissions").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateRoles(id string, roles []domain.Role) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Pastikan user ada
		var user domain.User
		if err := tx.Select("id").First(&user, "id = ?", id).Error; err != nil {
			return err
		}

		// Hapus relasi lama untuk user ini (filter juga model_type agar aman)
		if err := tx.
			Where("model_uuid = ? AND model_type = ?", id, domain.UserModelType).
			Delete(&domain.ModelHasRole{}).Error; err != nil {
			return err
		}

		// Jika tidak ada roles baru, selesai
		if len(roles) == 0 {
			return nil
		}

		// Insert relasi baru
		rows := make([]domain.ModelHasRole, 0, len(roles))
		for _, role := range roles {
			rows = append(rows, domain.ModelHasRole{
				RoleID:    role.ID,
				ModelUUID: id,
				ModelType: domain.UserModelType,
			})
		}
		if err := tx.Create(&rows).Error; err != nil {
			return err
		}

		return nil
	})
}
