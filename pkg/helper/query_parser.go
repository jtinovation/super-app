// helper/query_parser.go

package helper

import (
	"encoding/json"
	"jti-super-app-go/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParsePaginationQuery(c *gin.Context) (*dto.QueryParams, bool) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	search := c.DefaultQuery("search", "")
	sort := c.DefaultQuery("sort", "")
	order := c.DefaultQuery("order", "asc")
	filterStr := c.DefaultQuery("filter", "{}")

	errorsMap := map[string]string{}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		errorsMap["page"] = "Page must be a valid positive integer"
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		errorsMap["per_page"] = "Per page must be a valid positive integer"
	}

	if order != "asc" && order != "desc" {
		errorsMap["order"] = "Order must be either 'asc' or 'desc'"
	}

	var filter map[string]string
	if err := json.Unmarshal([]byte(filterStr), &filter); err != nil && filterStr != "{}" && filterStr != "" {
		errorsMap["filter"] = "Filter must be a valid JSON object"
	}

	if len(errorsMap) > 0 {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Message: "Validation errors occurred",
			Errors:  errorsMap,
		})
		return nil, false
	}

	return &dto.QueryParams{
		Page:    page,
		PerPage: perPage,
		Search:  SanitizeInput(search),
		Sort:    SanitizeInput(sort),
		Order:   order,
		Filter:  filter,
	}, true
}
