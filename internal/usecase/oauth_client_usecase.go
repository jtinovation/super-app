package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type OauthClientUseCase interface {
	FindAll(params dto.QueryParams) (*[]domain.OauthClient, int64, error)
	FindByID(id string) (*domain.OauthClient, error)
	Create(dto *dto.StoreOauthClientDTO) (*domain.OauthClient, error)
	Update(id string, dto *dto.UpdateOauthClientDTO) (*domain.OauthClient, error)
	Delete(id string) error
}

type oauthClientUseCase struct {
	repo domain.OauthClientRepository
}

func NewOauthClientUseCase(repo domain.OauthClientRepository) OauthClientUseCase {
	return &oauthClientUseCase{repo: repo}
}

func (u *oauthClientUseCase) FindAll(params dto.QueryParams) (*[]domain.OauthClient, int64, error) {
	return u.repo.FindAll(params)
}

func (u *oauthClientUseCase) FindByID(id string) (*domain.OauthClient, error) {
	return u.repo.FindByID(id)
}

func (u *oauthClientUseCase) Create(dto *dto.StoreOauthClientDTO) (*domain.OauthClient, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(dto.Secret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	client := &domain.OauthClient{
		ID:       uuid.NewString(),
		Name:     dto.Name,
		Secret:   string(hashedSecret),
		Redirect: dto.Redirect,
	}

	err = u.repo.Create(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (u *oauthClientUseCase) Update(id string, dto *dto.UpdateOauthClientDTO) (*domain.OauthClient, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(dto.Secret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	client := &domain.OauthClient{
		Name:     dto.Name,
		Secret:   string(hashedSecret),
		Redirect: dto.Redirect,
	}

	err = u.repo.Update(id, client)
	if err != nil {
		return nil, err
	}

	updatedClient, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return updatedClient, nil
}

func (u *oauthClientUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}
