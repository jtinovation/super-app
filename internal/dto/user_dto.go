package dto

import "time"

type UserResource struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Status      string     `json:"status"`
	Gender      *string    `json:"gender,omitempty"`
	Religion    *string    `json:"religion,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	BirthPlace  *string    `json:"birth_place,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	Address     *string    `json:"address,omitempty"`
	Nationality *string    `json:"nationality,omitempty"`
	Avatar      string     `json:"avatar"`
}
