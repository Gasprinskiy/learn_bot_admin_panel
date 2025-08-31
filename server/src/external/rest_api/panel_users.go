package rest_api

import (
	"context"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/external/rest_api/middleware"
	"learn_bot_admin_panel/internal/entity/app_jwt"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/tools/gin_gen"
	"learn_bot_admin_panel/tools/logger"
	"learn_bot_admin_panel/uimport"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PanelUsersHandler struct {
	ui         *uimport.UsecaseImport
	router     *gin.RouterGroup
	config     *config.Config
	log        *logger.Logger
	middleware *middleware.AuthMiddleware
	sm         transaction.SessionManager
}

func NewPanelUsersHandler(
	ui *uimport.UsecaseImport,
	router *gin.RouterGroup,
	config *config.Config,
	log *logger.Logger,
	middleware *middleware.AuthMiddleware,
	sm transaction.SessionManager,
) {
	handler := PanelUsersHandler{
		ui,
		router,
		config,
		log,
		middleware,
		sm,
	}

	group := handler.router.Group("/panel_users")

	{
		group.GET(
			"",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull}),
			handler.FindUsersList,
		)

		group.POST(
			"/create",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull}),
			handler.CreateProfile,
		)

		group.PATCH(
			"/redact",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull}),
			handler.RedactProfile,
		)
	}
}

func (h *PanelUsersHandler) FindUsersList(gctx *gin.Context) {
	jwtClaims, exists := gctx.Get(middleware.UserDataKey)
	if !exists {
		gin_gen.HandleError(gctx, global.ErrInternalError)
		return
	}

	claimsData, ok := jwtClaims.(app_jwt.Claims)
	if !ok {
		gin_gen.HandleError(gctx, global.ErrInternalError)
		return
	}

	data, err := transaction.RunInTx(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) ([]profile.User, error) {
			return h.ui.Usecase.Profile.FindProfileUsers(ctx, claimsData.UserID)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, data)
}

func (h *PanelUsersHandler) CreateProfile(gctx *gin.Context) {
	var param profile.CreateProfileParam

	if err := gctx.ShouldBindQuery(&param); err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	_, err := transaction.RunInTxCommit(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) (int64, error) {
			return h.ui.Usecase.Profile.CreateProfile(ctx, param)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *PanelUsersHandler) RedactProfile(gctx *gin.Context) {
	var param profile.RedactProfileParam

	if err := gctx.ShouldBindQuery(&param); err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	err := transaction.RunInTxExec(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) error {
			return h.ui.Usecase.Profile.RedactProfile(ctx, param)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, gin.H{"success": true})
}
