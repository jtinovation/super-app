package middleware

import (
	"context"
	"errors"
	"jti-super-app-go/config"
	"jti-super-app-go/internal/service"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helper.ErrorResponse(c, http.StatusUnauthorized, "Missing Authorization header", errors.New("unauthorized"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helper.ErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header format", errors.New("unauthorized"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		val, err := config.Rdb.Get(context.Background(), tokenString).Result()
		if err == nil && val == "blacklisted" {
			helper.ErrorResponse(c, http.StatusUnauthorized, "Token has been invalidated", errors.New("unauthorized"))
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Set("permissions", claims.Permissions)
		c.Set("token", tokenString)
		c.Next()
	}
}

func sliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Authorize checks if the user has the required roles or permissions.
// `requirements` is a pipe-separated string like "role:admin|super-admin" or "permission:edit-major"
func Authorize(requirements string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, _ := c.Get("roles")
		userPermissions, _ := c.Get("permissions")

		roles, okRoles := userRoles.([]string)
		if !okRoles {
			helper.ErrorResponse(c, http.StatusForbidden, "You are not authorized to perform this action", errors.New("unauthorized"))
			c.Abort()
			return
		}
		
		permissions, okPermissions := userPermissions.([]string)
		if !okPermissions {
			helper.ErrorResponse(c, http.StatusForbidden, "You are not authorized to perform this action", errors.New("unauthorized"))
			c.Abort()
			return
		}
		
		requiredItems := strings.Split(requirements, "|")
		
		isAuthorized := false
		for _, item := range requiredItems {
			parts := strings.SplitN(item, ":", 2)
			if len(parts) != 2 {
				continue
			}
			reqType := parts[0]
			reqValue := parts[1]
			
			if reqType == "role" && sliceContains(roles, reqValue) {
				isAuthorized = true
				break
			}
			if reqType == "permission" && sliceContains(permissions, reqValue) {
				isAuthorized = true
				break
			}
		}
		
		if !isAuthorized {
			helper.ErrorResponse(c, http.StatusForbidden, "You do not have the required permissions", errors.New("forbidden"))
			c.Abort()
			return
		}
		
		c.Next()
	}
}