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
	"learn_bot_admin_panel/tools/gin_gen"
	"learn_bot_admin_panel/tools/logger"
	"learn_bot_admin_panel/uimport"
	"net/http"
	"strconv"

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

		group.GET(
			"/excel_file",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull, profile.AccessRightManager}),
			handler.PrintBotUsers,
		)
	}
}

func (h *BotUsersHandler) GetBotUsers(gctx *gin.Context) {
	var param bot_users.FindBotRegisteredUsersQuertParseParam

	if err := gctx.ShouldBindQuery(&param); err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	data, err := transaction.RunInTx(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) (global.CommotListSearchResponse[bot_users.BotUserProfile], error) {
			return h.ui.Usecase.BotUsers.FindRegisteredUsers(ctx, param.InnerParam())
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, data)
}

func (h *BotUsersHandler) PrintBotUsers(gctx *gin.Context) {
	var param bot_users.FindBotRegisteredUsersQuertParseParam

	if err := gctx.ShouldBindQuery(&param); err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	data, err := transaction.RunInTx(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) ([]byte, error) {
			return h.ui.Usecase.BotUsers.PrintFindRegisteredUsers(ctx, param.InnerParam())
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.Header("Content-Description", "File Transfer")
	gctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "bot_users.xlsx"))
	gctx.Header("Content-Length", strconv.Itoa(len(data)))

	gctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
