package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CurrentCompanyID(c *gin.Context) uint {
	session := sessions.Default(c)

	value := session.Get("company_id")
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case int64:
		return uint(v)
	case float64:
		return uint(v)
	default:
		return 0
	}
}
