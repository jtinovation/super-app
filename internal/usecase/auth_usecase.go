package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"jti-super-app-go/config"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/service"
	"jti-super-app-go/pkg/constants"
	"jti-super-app-go/pkg/helper"
	"slices"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(req dto.LoginRequestDTO) (*dto.LoginResponseDTO, error)
	LoginWithGoogle(userInfo *service.GoogleUserInfo) (*dto.LoginResponseDTO, error)
	Logout(tokenString string) error
	VerifyEmail(token string) error
	ResendVerificationEmail(email string) error
	ForgotPassword(req dto.ForgotPasswordRequestDTO) error
	ResetPassword(req dto.ResetPasswordRequestDTO) error
	Me(userID string) (*dto.UserDetailInfoDTO, error)
}

type authUseCase struct {
	authRepo      domain.AuthRepository
	userRepo      domain.UserRepository
	employeeRepo  domain.EmployeeRepository
	studentRepo   domain.StudentRepository
	passResetRepo domain.PasswordResetRepository
	jwtService    service.JWTService
	emailService  service.EmailService
}

func NewAuthUseCase(authRepo domain.AuthRepository, userRepo domain.UserRepository, employeeRepo domain.EmployeeRepository, studentRepo domain.StudentRepository, passResetRepo domain.PasswordResetRepository, jwtService service.JWTService, emailService service.EmailService) AuthUseCase {
	return &authUseCase{
		authRepo:      authRepo,
		userRepo:      userRepo,
		employeeRepo:  employeeRepo,
		studentRepo:   studentRepo,
		passResetRepo: passResetRepo,
		jwtService:    jwtService,
		emailService:  emailService,
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

	if user.EmailVerifiedAt == nil {
		go uc.ResendVerificationEmail(user.Email)

		return nil, errors.New("please verify your email address, a new verification link has been sent")
	}

	var roleNames []string
	var permissionNames []string
	permissionSet := make(map[string]struct{})

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
			ID:               user.ID,
			Name:             user.Name,
			Email:            user.Email,
			IsChangePassword: user.IsChangePassword,
			Roles:            roleNames,
			Permissions:      permissionNames,
		},
	}, nil
}

func (uc *authUseCase) LoginWithGoogle(userInfo *service.GoogleUserInfo) (*dto.LoginResponseDTO, error) {
	user, err := uc.authRepo.FindByEmail(userInfo.Email)
	if err != nil {
		return nil, errors.New("email not registered")
	}

	var roleNames []string
	var permissionNames []string
	permissionSet := make(map[string]struct{})

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
			ID:               user.ID,
			Name:             user.Name,
			Email:            user.Email,
			IsChangePassword: user.IsChangePassword,
			Roles:            roleNames,
			Permissions:      permissionNames,
		},
	}, nil
}

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

func (uc *authUseCase) ForgotPassword(req dto.ForgotPasswordRequestDTO) error {
	user, err := uc.authRepo.FindByEmail(req.Email)
	if err != nil {
		return nil
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)

	pr := &domain.PasswordReset{
		Email:     user.Email,
		Token:     token,
		CreatedAt: time.Now(),
	}

	if _, err := uc.passResetRepo.Create(pr); err != nil {
		return err
	}

	resetLink := config.AppConfig.FrontendURL + "/reset-password?token=" + token
	subject := "Reset Your Password"
	templatePath := "templates/auth/password_reset.html"
	data := dto.EmailTemplateAuthDataDto{
		Name:    user.Name,
		Link:    resetLink,
		LogoURL: helper.GetUrlFile(constants.EMPLOYEE_PATH, constants.DEFAULT_AVATAR),
	}

	if err := uc.emailService.SendEmailWithTemplate(user.Email, subject, templatePath, data); err != nil {
		return err
	}

	return nil
}

func (uc *authUseCase) ResetPassword(req dto.ResetPasswordRequestDTO) error {
	pr, err := uc.passResetRepo.FindByTokenAndEmail(req.Token, req.Email)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Since(pr.CreatedAt) > time.Hour*1 {
		uc.passResetRepo.Delete(req.Token)
		return errors.New("invalid or expired token")
	}

	user, err := uc.authRepo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.IsChangePassword = true

	if _, err := uc.userRepo.Update(user.ID, user); err != nil {
		return err
	}

	uc.passResetRepo.Delete(req.Token)

	return nil
}

// Implementasi untuk Verifikasi Email
func (uc *authUseCase) VerifyEmail(token string) error {
	// Logika untuk verifikasi email. Anda bisa menggunakan Redis untuk menyimpan token verifikasi.
	// Contoh:
	email, err := config.Rdb.Get(context.Background(), "verify_email:"+token).Result()
	if err != nil {
		return errors.New("invalid or expired token")
	}

	user, err := uc.authRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	if _, err := uc.userRepo.Update(user.ID, user); err != nil {
		return err
	}

	// Hapus token dari Redis
	config.Rdb.Del(context.Background(), "verify_email:"+token)
	return nil
}

func (uc *authUseCase) ResendVerificationEmail(email string) error {
	user, err := uc.authRepo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	if user.EmailVerifiedAt != nil {
		return errors.New("email already verified")
	}

	// Buat token baru
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)

	// Simpan ke Redis dengan TTL (misal: 15 menit)
	err = config.Rdb.Set(context.Background(), "verify_email:"+token, user.Email, 15*time.Minute).Err()
	if err != nil {
		return err
	}

	verificationLink := "http://your-backend-app.com/api/v1/auth/email/verify/" + token
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`
		<h1>Welcome to JTI Super App!</h1>
		<p>Please verify your email address by clicking the link below:</p>
		<a href="%s">Verify Email</a>
		<p>This link will expire in 15 minutes.</p>
	`, verificationLink)
	if err := uc.emailService.SendEmail(user.Email, subject, body); err != nil {
		return err
	}

	// TODO: Kirim email verifikasi ke pengguna
	// Contoh: `http://frontend.com/verify-email?token=` + token
	fmt.Println("Verification Email Token:", token) // Untuk debugging

	return nil
}

func (uc *authUseCase) Me(userID string) (*dto.UserDetailInfoDTO, error) {
	cacheKey := "user_info:" + userID
	var userInfo dto.UserDetailInfoDTO
	var employee *domain.Employee
	var student *domain.Student
	studentSemesters := &[]dto.StudentSemesterDTO{}

	cached, err := config.Rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cached != "" {
		if err := json.Unmarshal([]byte(cached), &userInfo); err == nil {
			return &userInfo, nil
		}
	}

	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	var roleNames []string
	var permissionNames []string
	permissionSet := make(map[string]struct{})

	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
		for _, perm := range role.Permissions {
			permissionSet[perm.Name] = struct{}{}
		}
	}

	for perm := range permissionSet {
		permissionNames = append(permissionNames, perm)
	}

	if slices.Contains(roleNames, "student") {
		employee = &domain.Employee{}
		stud, err := uc.studentRepo.FindByUserID(userID)
		if err == nil {
			student = stud
			if student != nil && len(student.StudentSemesters) > 0 {
				var semesters []dto.StudentSemesterDTO
				for _, ss := range student.StudentSemesters {
					semesters = append(semesters, dto.StudentSemesterDTO{
						ID:         ss.ID,
						SemesterID: ss.SemesterID,
						StudentID:  ss.StudentID,
						Class:      ss.Class,
						IsActive:   ss.IsActive,
					})
				}
				studentSemesters = &semesters
			}
		} else {
			student = &domain.Student{}
		}
	} else {
		student = &domain.Student{}
		emp, err := uc.employeeRepo.FindByUserID(userID)
		if err == nil {
			employee = emp
		} else {
			employee = &domain.Employee{}
		}
	}

	userInfo = dto.UserDetailInfoDTO{
		ID:               user.ID,
		Name:             user.Name,
		Email:            user.Email,
		EmailVerifiedAt:  user.EmailVerifiedAt,
		Status:           user.Status,
		Gender:           user.Gender,
		Religion:         user.Religion,
		BirthPlace:       user.BirthPlace,
		BirthDate:        user.BirthDate,
		PhoneNumber:      user.PhoneNumber,
		Address:          user.Address,
		Nationality:      user.Nationality,
		ImgPath:          user.ImgPath,
		ImgName:          user.ImgName,
		IsChangePassword: user.IsChangePassword,
		Roles:            roleNames,
		Permissions:      permissionNames,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		DeletedAt:        user.DeletedAt,
		EmployeeDetail: &dto.EmployeeDetailInfoDTO{
			ID:               employee.ID,
			NIP:              employee.Nip,
			MajorID:          employee.MajorID,
			MajorName:        &employee.Major.Name,
			StudyProgramID:   employee.StudyProgramID,
			StudyProgramName: &employee.StudyProgram.Name,
			Position:         employee.Position,
			CreatedAt:        &employee.CreatedAt,
			UpdatedAt:        &employee.UpdatedAt,
		},
		StudentDetail: &dto.StudentDetailInfoDTO{
			ID:               student.ID,
			NIM:              student.NIM,
			Generation:       student.Generation,
			MajorID:          &student.StudyProgram.MajorID,
			MajorName:        &student.StudyProgram.Major.Name,
			StudyProgramID:   &student.StudyProgram.ID,
			StudyProgramName: &student.StudyProgram.Name,
			StudentSemesters: studentSemesters,
			CreatedAt:        &student.CreatedAt,
			UpdatedAt:        &student.UpdatedAt,
		},
	}

	// Marshal userInfo to JSON and store in Redis with TTL (e.g., 10 minutes)
	if jsonBytes, err := json.Marshal(userInfo); err == nil {
		_ = config.Rdb.Set(context.Background(), cacheKey, jsonBytes, 10*time.Minute).Err()
	}

	return &userInfo, nil
}
