package rest_api

import (
	"context"
	"fmt"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/external/rest_api/middleware"
	"learn_bot_admin_panel/internal/entity/bot_users"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/tools/dump"
	"learn_bot_admin_panel/tools/gin_gen"
	"learn_bot_admin_panel/tools/logger"
	"learn_bot_admin_panel/uimport"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BotUsersHandler struct {
	ui         *uimport.UsecaseImport
	router     *gin.RouterGroup
	config     *config.Config
	log        *logger.Logger
	middleware *middleware.AuthMiddleware
	sm         transaction.SessionManager
}

func NewBotUsersHandler(
	ui *uimport.UsecaseImport,
	router *gin.RouterGroup,
	config *config.Config,
	log *logger.Logger,
	middleware *middleware.AuthMiddleware,
	sm transaction.SessionManager,
) {
	handler := BotUsersHandler{
		ui,
		router,
		config,
		log,
		middleware,
		sm,
	}

	group := handler.router.Group("/bot_users")

	{
		group.GET(
			"",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull, profile.AccessRightManager}),
			handler.GetBotUsers,
		)
	}
}

func (h *BotUsersHandler) GetBotUsers(gctx *gin.Context) {
	var param bot_users.FindBotRegisteredUsersQuertParseParam

	if err := gctx.ShouldBindQuery(&param); err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	innerParam := param.InnerParam()

	fmt.Println("innerParam: ", dump.Struct(innerParam))

	data, err := transaction.RunInTx(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) (global.CommotListSearchResponse[bot_users.BotUserProfile], error) {
			return h.ui.Usecase.BotUsers.FindRegisteredUsers(ctx, innerParam)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, data)
}
