package helper

import (
	"errors"
	"jti-super-app-go/internal/dto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationErrorJSON(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errorsList := make(map[string]string)
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			tag := fe.Tag()
			var msg string

			switch tag {
			case "required":
				msg = field + " is required"
			case "max":
				msg = field + " must not exceed " + fe.Param() + " characters"
			case "min":
				msg = field + " must be at least " + fe.Param() + " characters"
			case "email":
				msg = field + " must be a valid email address"
			case "uuid":
				msg = field + " must be a valid UUID"
			case "gte":
				msg = field + " must be greater than or equal to " + fe.Param()
			case "lte":
				msg = field + " must be less than or equal to " + fe.Param()
			case "oneof":
				msg = field + " must be one of the following values: " + strings.ReplaceAll(fe.Param(), " ", ", ")
			case "url":
				msg = field + " must be a valid URL"
			case "alpha":
				msg = field + " must contain only alphabetic characters"
			case "numeric":
				msg = field + " must contain only numeric characters"
			case "alphanum":
				msg = field + " must contain only alphanumeric characters"
			case "datetime":
				msg = field + " must be a valid datetime in the format " + fe.Param()
			default:
				msg = field + " is invalid"
			}

			errorsList[field] = msg
		}

		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Message: "Invalid request data",
			Errors:  errorsList,
		})
		return
	}

	c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Message: "Invalid request payload",
		Errors:  map[string]string{"details": err.Error()},
	})
}
