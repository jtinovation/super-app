package handler

import (
	"encoding/base64"
	"encoding/json"
	"jti-super-app-go/config"
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type OauthHandler struct {
	useCase      usecase.OauthClientUseCase
	oauthUseCase usecase.OauthUsecase
	authUseCase  usecase.AuthUseCase
}

func NewOauthHandler(uc usecase.OauthClientUseCase, oc usecase.OauthUsecase, ac usecase.AuthUseCase) *OauthHandler {
	return &OauthHandler{useCase: uc, oauthUseCase: oc, authUseCase: ac}
}

func (h *OauthHandler) Authorize(c *gin.Context) {
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	responseType := c.Query("response_type")

	if clientID == "" || redirectURI == "" || responseType != "code" {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			"Invalid request parameters")
		return
	}

	client, err := h.useCase.FindByID(clientID)
	if err != nil || client.Redirect != redirectURI {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			"Invalid client or redirect URI")
		return
	}

	userId, ok := helper.GetSSO(c) // baca cookie sso_session
	if !ok {
		config.Rdb.Set(c, "return_to", c.Request.URL.RequestURI(), 0)

		u := url.URL{Path: "/login"}
		q := u.Query()
		q.Set("return_to", c.Request.URL.RequestURI())
		u.RawQuery = q.Encode()
		c.Redirect(http.StatusSeeOther, u.String())
		return
	}

	dataCookie, err := config.Rdb.Get(c, "sso:"+userId).Result()
	if err != nil {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			"SSO session expired or invalid")
		return
	}

	var user dto.LoginResponseDTO
	err = json.Unmarshal([]byte(dataCookie), &user)
	if err != nil {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			"Failed to parse user data")
		return
	}

	data, err := h.oauthUseCase.Authorize(clientID, redirectURI, responseType, &user)
	if err != nil {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			"Failed to generate authorization code")
		return
	}

	u, err := url.Parse(redirectURI)
	if err != nil {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			"Invalid redirect URI")
		return
	}

	q := u.Query()
	q.Set("code", data.Code)
	u.RawQuery = q.Encode()
	c.Redirect(http.StatusSeeOther, u.String())
}

func (h *OauthHandler) Token(c *gin.Context) {
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	redirectURI := c.PostForm("redirect_uri")
	code := c.PostForm("code")

	if clientID == "" || clientSecret == "" || redirectURI == "" || code == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request parameters", nil)
		return
	}

	client, err := h.useCase.FindByID(clientID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid client or client secret", err)
		return
	}

	if client.Redirect != redirectURI {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid redirect URI", nil)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(client.Secret), []byte(clientSecret))
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Invalid client", err)
		return
	}

	val, err := config.Rdb.Get(c, code).Result()
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid authorization code", err)
		return
	}

	var ac dto.StoreOauthCodeDTO
	if err := json.Unmarshal([]byte(val), &ac); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to parse authorization code data", err)
		return
	}

	if time.Now().After(ac.ExpiresAt) || ac.ClientID != client.ID {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid grant", nil)
		return
	}

	if ac.RedirectURI != redirectURI {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid redirect URI", nil)
		return
	}

	if ac.Code != code {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid authorization code", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Token generated successfully", ac.UserSub)
}

func (h *OauthHandler) LoginPage(c *gin.Context) {
	returnTo := c.Query("return_to")
	token, _ := c.Get("csrf_token")
	var errMsg string
	if enc := c.Query("error"); enc != "" {
		if b, decErr := base64.RawURLEncoding.DecodeString(enc); decErr == nil {
			errMsg = string(b)
		}
	}

	if _, err := c.Cookie(helper.CookieName); err == nil {
		c.Redirect(http.StatusSeeOther, returnTo)
		return
	}

	c.HTML(http.StatusOK, "auth/login.tmpl", gin.H{
		"return_to":  returnTo,
		"csrf_token": token,
		"error":      errMsg,
	})
}

func (h *OauthHandler) LoginPost(c *gin.Context) {
	var form dto.LoginRequestFormDTO
	if err := c.ShouldBind(&form); err != nil {
		helper.RedirectBackToLogin(c, "/login", form.ReturnTo, "",
			"Invalid form data")
		return
	}

	if !helper.ValidateCSRF(c) {
		helper.RedirectBackToLogin(c, "/login", form.ReturnTo, "",
			"Invalid CSRF token")
		return
	}

	jsonReq := dto.LoginRequestDTO{
		Email:    form.Email,
		Password: form.Password,
	}

	user, err := h.authUseCase.Login(jsonReq)
	if err != nil {
		helper.RedirectBackToLogin(c, "/login", form.ReturnTo, "",
			"Login failed: "+err.Error())
		return
	}

	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	helper.SetSSO(c, user, config.AppConfig.JWTExpirationHours*3600) // simpan cookie SSO selama JWTExpirationHours

	c.Redirect(http.StatusSeeOther, form.ReturnTo)
}

func (h *OauthHandler) LoginCallback(c *gin.Context) {
	token := c.Query("token")
	userData := c.Query("user")
	errorMsg := c.Query("error")

	if token == "" || userData == "" {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			"Missing token or user data")
		return
	}

	if errorMsg != "" {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			errorMsg)
		return
	}

	var user dto.LoginResponseDTO
	decodedUserData, err := base64.RawURLEncoding.DecodeString(userData)
	if err != nil {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			"Failed to decode user data")
		return
	}

	err = json.Unmarshal(decodedUserData, &user.User)
	if err != nil {
		helper.ClearSSO(c)
		helper.RedirectBackToLogin(c, "/login", c.Request.URL.RequestURI(), "",
			"Failed to parse user data")
		return
	}

	user.Token = token

	helper.SetSSO(c, &user, config.AppConfig.JWTExpirationHours*3600) // simpan cookie SSO selama JWTExpirationHours
	return_to, _ := config.Rdb.Get(c, "return_to").Result()
	config.Rdb.Del(c, "return_to")
	c.Redirect(http.StatusSeeOther, return_to)
}

func (h *OauthHandler) IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/index.tmpl", gin.H{})
}

func (h *OauthHandler) Logout(c *gin.Context) {
	helper.ClearSSO(c)
	redirectTo := c.Query("redirect")
	c.Redirect(http.StatusSeeOther, redirectTo)
}
