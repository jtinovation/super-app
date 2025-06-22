package usecase

import (
	"context"
	"errors"
	"jti-super-app-go/config"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/service"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(req dto.LoginRequestDTO) (*dto.LoginResponseDTO, error)
	Logout(tokenString string) error
}

type authUseCase struct {
	authRepo   domain.AuthRepository
	jwtService service.JWTService
}

func NewAuthUseCase(authRepo domain.AuthRepository, jwtService service.JWTService) AuthUseCase {
	return &authUseCase{
		authRepo:   authRepo,
		jwtService: jwtService,
	}
}

func (uc *authUseCase) Login(req dto.LoginRequestDTO) (*dto.LoginResponseDTO, error) {
	user, err := uc.authRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Extract role and permission names
	var roleNames []string
	var permissionNames []string
	permissionSet := make(map[string]struct{}) // Use a map to store unique permissions

	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
		for _, perm := range role.Permissions {
			permissionSet[perm.Name] = struct{}{}
		}
	}
    
	for perm := range permissionSet {
		permissionNames = append(permissionNames, perm)
	}

	token, err := uc.jwtService.GenerateToken(user.ID, roleNames, permissionNames)
	if err != nil {
		return nil, errors.New("could not generate token")
	}

	return &dto.LoginResponseDTO{
		Token: token,
		User: dto.UserLoginInfo{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Roles:     roleNames,
			Permissions: permissionNames,
		},
	}, nil
}

// the logout method
func (uc *authUseCase) Logout(tokenString string) error {
	claims, err := uc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil
	}

	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		return nil
	}

	err = config.Rdb.Set(context.Background(), tokenString, "blacklisted", remaining).Err()
	if err != nil {
		return errors.New("failed to blacklist token")
	}

	return nil
}