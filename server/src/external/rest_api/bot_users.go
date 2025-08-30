package rest_api

import (
	"context"
	"fmt"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/external/rest_api/middleware"
	"learn_bot_admin_panel/internal/entity/app_jwt"
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
			"/registered",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull, profile.AccessRightManager}),
			handler.GetBotUsers,
		)

		group.GET(
			"/unregistered",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull, profile.AccessRightManager}),
			handler.GetUnregisteredBotUsers,
		)

		group.GET(
			"/registered/excel_file",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull, profile.AccessRightManager}),
			handler.GetBotUsersExcelFile,
		)

		group.GET(
			"/unregistered/excel_file",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull, profile.AccessRightManager}),
			handler.GetBotUnregisteredUsersExcelFile,
		)

		group.GET(
			"/:user_id",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull, profile.AccessRightManager}),
			handler.GetUserByID,
		)

		group.GET(
			"/subscr_types",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull, profile.AccessRightManager}),
			handler.GetBotSubscriptionTypes,
		)

		group.POST(
			"/purchase/:user_id",
			handler.middleware.CheckAccesToken(),
			handler.middleware.CheckAccessRight([]profile.AccessRight{profile.AccessRightFull, profile.AccessRightManager}),
			handler.PostPurchase,
		)
	}
}

func (h *BotUsersHandler) GetBotUsers(gctx *gin.Context) {
	var param bot_users.FindBotRegisteredUsersQueryParseParam

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

func (h *BotUsersHandler) GetUnregisteredBotUsers(gctx *gin.Context) {
	var param bot_users.FindBotUnregisteredUsersQueryParseParam

	if err := gctx.ShouldBindQuery(&param); err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	data, err := transaction.RunInTx(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) (global.CommotListSearchResponse[bot_users.BotUnregistredUserProfile], error) {
			return h.ui.Usecase.BotUsers.FindUnregisteredUsers(ctx, param.InnerParam())
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, data)
}

func (h *BotUsersHandler) GetBotUsersExcelFile(gctx *gin.Context) {
	var param bot_users.FindBotRegisteredUsersQueryParseParam

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

func (h *BotUsersHandler) GetBotUnregisteredUsersExcelFile(gctx *gin.Context) {
	var param bot_users.FindBotUnregisteredUsersQueryParseParam

	if err := gctx.ShouldBindQuery(&param); err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	data, err := transaction.RunInTx(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) ([]byte, error) {
			return h.ui.Usecase.BotUsers.PrintFindUnregisteredUsers(ctx, param.InnerParam())
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

func (h *BotUsersHandler) GetUserByID(gctx *gin.Context) {
	userID, err := strconv.Atoi(gctx.Param("user_id"))
	if err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	data, err := transaction.RunInTx(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) (bot_users.BotUserDetailData, error) {
			return h.ui.Usecase.BotUsers.FindUserByID(ctx, userID)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, data)
}

func (h *BotUsersHandler) PostPurchase(gctx *gin.Context) {
	userID, err := strconv.Atoi(gctx.Param("user_id"))
	if err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	serviceID, err := strconv.Atoi(gctx.PostForm("sub_id"))
	if err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	file, header, err := gctx.Request.FormFile("receipt")
	if err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}
	defer file.Close()

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

	param := bot_users.NewPurchaseSubscriptionParam(
		userID,
		claimsData.UserID,
		serviceID,
		file,
		*header,
	)

	err = transaction.RunInTxExec(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) error {
			return h.ui.BotUsers.PurchaseSubscription(ctx, param)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *BotUsersHandler) GetBotSubscriptionTypes(gctx *gin.Context) {
	data, err := transaction.RunInTx(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) ([]bot_users.BotSubscriptionType, error) {
			return h.ui.Usecase.BotUsers.LoadAllBotSubscriptionTypes(ctx)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, data)
}
