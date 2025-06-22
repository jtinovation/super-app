package helper

import (
	"errors"
	"jti-super-app-go/internal/dto"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, dto.SingleResponse{
		Message: message,
		Data:    data,
	})
}

func PaginatedSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, meta *dto.Meta) {
	c.JSON(statusCode, dto.PaginatedResponse{
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	if err == nil {
		err = gin.Error{Err: errors.New("unknown error")}
	}

	c.AbortWithStatusJSON(statusCode, dto.ErrorResponse{
		Message: message,
		Errors:  map[string]string{"error": err.Error()},
	})
}