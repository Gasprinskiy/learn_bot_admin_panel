package restapi

import (
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/event_bus"
	"learn_bot_admin_panel/uimport"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	ui       *uimport.UsecaseImport
	gin      *gin.Engine
	config   *config.Config
	uuidChan event_bus.UUIDChanel
}

func NewProfileHandler(
	ui *uimport.UsecaseImport,
	gin *gin.Engine,
	config *config.Config,
	uuidChan event_bus.UUIDChanel,
) {
	handler := ProfileHandler{
		ui,
		gin,
		config,
		uuidChan,
	}

	handler.gin.Group("/auth")

	{
		handler.gin.GET(
			"/temp_id",
			GetTemporaryID,
		)

		handler.gin.GET(
			"/listen",
		)
	}
}

func GetTemporaryID(ctx *gin.Context) {
	// uuID := uuid.NewString()

	// ctx.JSON(http.StatusCreated)
}
