package utils

import (
	"findai/src/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Paginate() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		if page < 1 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		if limit > 100 {
			limit = 100
		}

		offset := (page - 1) * limit

		// Collect filters from query parameters
		var filters []database.Filter
		for key, values := range c.Request.URL.Query() {
			if key != "page" && key != "limit" {
				for _, value := range values {
					filters = append(filters, database.Filter{Key: key, Value: value})
				}
			}
		}

		pagination := database.Paginate{
			Limit:   limit,
			Offset:  offset,
			Filters: filters,
		}

		c.Set("paginate", pagination)
		c.Set("page", page)
		c.Set("limit", limit)
		c.Next()
	}
}
