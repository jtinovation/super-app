package usecase

import (
	"context"
	"encoding/json"
	"jti-super-app-go/config"
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/pkg/helper"
	"time"
)

type OauthUsecase interface {
	Authorize(clientID, redirectURI, responseType string, user *dto.LoginResponseDTO) (dto.StoreOauthCodeDTO, error)
}

type oauthUsecase struct {
	// dependencies can be added here
}

func NewOauthUsecase() OauthUsecase {
	return &oauthUsecase{}
}

func (uc *oauthUsecase) Authorize(clientID, redirectURI, responseType string, user *dto.LoginResponseDTO) (dto.StoreOauthCodeDTO, error) {
	code := helper.GenCode()

	data := dto.StoreOauthCodeDTO{
		Code:        code,
		ClientID:    clientID,
		UserSub:     *user,
		RedirectURI: redirectURI,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return dto.StoreOauthCodeDTO{}, err
	}

	storeToRedis := config.Rdb.Set(context.Background(), code, jsonData, 10*time.Minute)
	if storeToRedis.Err() != nil {
		return dto.StoreOauthCodeDTO{}, storeToRedis.Err()
	}

	return data, nil
}
