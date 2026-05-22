package pkg

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseInt64Param(c *gin.Context, paramName string) (int64, error) {
	paramStr := c.Param(paramName)
	id, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil {
		return 0, ErrIDInvalido
	}

	if id <= 0 {
		return 0, ErrIDInvalido
	}

	return id, nil
}
