package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/zzzgydi/zbyai/common"
	"github.com/zzzgydi/zbyai/model"
)

func GetUser(c *gin.Context) *model.User {
	user, _ := c.Get(common.CTX_CURRENT_USER)
	if user != nil {
		if u, ok := user.(*model.User); ok {
			return u
		}
	}
	return nil
}
