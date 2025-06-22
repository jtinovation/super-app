package domain

import (
	"time"
)

type User struct {
	ID                string     `gorm:"type:char(36);primaryKey"`
	Name              string     `gorm:"type:varchar(255);not null"`
	Email             string     `gorm:"type:varchar(255);unique;not null"`
	EmailVerifiedAt   *time.Time `gorm:"type:timestamp"`
	Password          string     `gorm:"type:varchar(255);not null"`
	RememberToken     *string    `gorm:"type:varchar(100)"`
	Status            string     `gorm:"type:enum('ACTIVE','INACTIVE');default:'ACTIVE';not null"`
	Gender            *string    `gorm:"type:enum('MALE','FEMALE')"`
	Religion          *string    `gorm:"type:enum('ISLAM','CHRISTIANITY','CATHOLIC','HINDUISM','BUDDHISM','CONFUCIANISM','OTHER')"`
	BirthPlace        *string    `gorm:"type:varchar(255)"`
	BirthDate         *time.Time `gorm:"type:date"`
	PhoneNumber       *string    `gorm:"type:varchar(255)"`
	Address           *string    `gorm:"type:varchar(255)"`
	Nationality       *string    `gorm:"type:varchar(255)"`
	ImgPath           *string    `gorm:"type:varchar(255)"`
	ImgName           *string    `gorm:"type:varchar(255)"`
	IsChangePassword  bool       `gorm:"type:tinyint(1);default:0;not null"`
	CreatedAt         *time.Time `gorm:"type:timestamp"`
	UpdatedAt         *time.Time `gorm:"type:timestamp"`
	DeletedAt         *time.Time `gorm:"type:timestamp"`
	Roles []Role `gorm:"many2many:model_has_roles;foreignKey:ID;joinForeignKey:model_uuid;References:ID;joinReferences:role_id"`
}

type Role struct {
	ID          string       `gorm:"type:char(36);primaryKey;column:uuid"`
	Name        string       `gorm:"type:varchar(255);unique;not null"`
	GuardName   string       `gorm:"type:varchar(255);default:'api'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Permissions []Permission `gorm:"many2many:role_has_permissions;foreignKey:ID;joinForeignKey:role_id;References:ID;joinReferences:permission_id"`
}

type Permission struct {
	ID        string `gorm:"type:char(36);primaryKey;column:uuid"`
	Name      string `gorm:"type:varchar(255);unique;not null"`
	GuardName string `gorm:"type:varchar(255);default:'api'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "m_user"
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}