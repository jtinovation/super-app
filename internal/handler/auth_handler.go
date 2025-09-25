package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/service"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/constants"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"net/url"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	useCase           usecase.AuthUseCase
	googleAuthService service.GoogleAuthService
}

func NewAuthHandler(uc usecase.AuthUseCase, gas service.GoogleAuthService) *AuthHandler {
	return &AuthHandler{useCase: uc, googleAuthService: gas}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	res, err := h.useCase.Login(req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login successful", res)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state := uuid.NewString()
	host := c.Request.Header.Get("Referer")
	fmt.Print("request header: ", c.Request.Header)
	fmt.Println("request host: ", host)

	c.SetCookie("oauth_state", state, 3600, "/api/v1/auth/google/callback", c.Request.URL.Hostname(), false, true)
	c.SetCookie("host", host, 3600, "/api/v1/auth/google/callback", c.Request.URL.Hostname(), false, true)

	url := h.googleAuthService.GenerateAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	stateFromCookie, err := c.Cookie("oauth_state")
	hostFromCookie, _ := c.Cookie("host")

	parsedHost := hostFromCookie
	if len(hostFromCookie) > 0 {
		u, err := url.Parse(hostFromCookie)
		if err == nil {
			parsedHost = u.Scheme + "://" + u.Host
		}
	}
	hostFromCookie = parsedHost

	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, hostFromCookie+constants.CALLBACK_FRONTEND+"?error="+base64.URLEncoding.EncodeToString([]byte("State cookie not found")))
		return
	}
	if c.Query("state") != stateFromCookie {
		c.Redirect(http.StatusTemporaryRedirect, hostFromCookie+constants.CALLBACK_FRONTEND+"?error="+base64.URLEncoding.EncodeToString([]byte("Invalid state parameter")))
		return
	}

	code := c.Query("code")
	userInfo, err := h.googleAuthService.GetUserInfo(code)
	if err != nil {
		encodeError := base64.URLEncoding.EncodeToString([]byte(err.Error()))
		c.Redirect(http.StatusTemporaryRedirect, hostFromCookie+constants.CALLBACK_FRONTEND+"?error="+encodeError)
		return
	}

	res, err := h.useCase.LoginWithGoogle(userInfo)
	if err != nil {
		encodeError := base64.URLEncoding.EncodeToString([]byte(err.Error()))
		c.Redirect(http.StatusTemporaryRedirect, hostFromCookie+constants.CALLBACK_FRONTEND+"?error="+encodeError)
		return
	}

	userInfoJSON, _ := json.Marshal(res.User)
	encodedUser := base64.URLEncoding.EncodeToString(userInfoJSON)
	c.Redirect(http.StatusPermanentRedirect, hostFromCookie+constants.CALLBACK_FRONTEND+"?token="+res.Token+"&user="+encodedUser)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Could not retrieve token from context", errors.New("token not found in context"))
		return
	}

	tokenString, ok := token.(string)
	if !ok {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Invalid token format in context", errors.New("invalid token format"))
		return
	}

	err := h.useCase.Logout(tokenString)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error(), errors.New("logout failed"))
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Successfully logged out", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	if err := h.useCase.ForgotPassword(req); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to send password reset link", err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Password reset link has been sent to your email", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	if err := h.useCase.ResetPassword(req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error(), err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Password has been successfully reset", nil)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	if err := h.useCase.VerifyEmail(token); err != nil {
		// Arahkan ke halaman error di frontend
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Arahkan ke halaman sukses di frontend
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	// Ambil email dari user yang sedang login (dari JWT)
	// Ini hanyalah contoh, sesuaikan dengan implementasi Anda
	// userID, _ := c.Get("user_id")
	// Anda perlu mengambil detail user dari ID ini untuk mendapatkan emailnya
	// Misal, Anda memiliki user usecase
	// user, err := h.userUseCase.FindByID(userID.(string))

	// Untuk sementara, kita ambil dari body
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	if err := h.useCase.ResendVerificationEmail(req.Email); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error(), err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Verification link has been sent to your email", nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Could not retrieve user ID from context", errors.New("user ID not found in context"))
		return
	}

	userInfo, err := h.useCase.Me(userID.(string))
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user information", err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "User information retrieved successfully", userInfo)
}
