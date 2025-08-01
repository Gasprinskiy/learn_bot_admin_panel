package gin_gen

import (
	"learn_bot_admin_panel/internal/entity/global"

	"github.com/gin-gonic/gin"
)

func HandleError(gctx *gin.Context, err error) {
	gctx.JSON(global.ErrStatusCodes[err], gin.H{"message": err.Error()})
}
