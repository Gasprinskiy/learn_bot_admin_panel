package rest_api

import (
	"context"
	"fmt"
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

type ProfileHandler struct {
	ui         *uimport.UsecaseImport
	router     *gin.RouterGroup
	config     *config.Config
	log        *logger.Logger
	middleware *middleware.AuthMiddleware
	sm         transaction.SessionManager
}

func NewProfileHandler(
	ui *uimport.UsecaseImport,
	router *gin.RouterGroup,
	config *config.Config,
	log *logger.Logger,
	middleware *middleware.AuthMiddleware,
	sm transaction.SessionManager,
) {
	handler := ProfileHandler{
		ui,
		router,
		config,
		log,
		middleware,
		sm,
	}

	group := handler.router.Group("/auth")

	{
		group.GET(
			"/temp_data",
			handler.GetAuthData,
		)

		group.GET(
			"/listen/:temp_id",
			handler.HandleAuthListen,
		)

		group.GET(
			"/tg_login/:temp_id",
			handler.HandleTgLogin,
		)

		group.POST(
			"/password_login",
			handler.HandleLogin,
		)

		group.GET(
			"/check",
			handler.HandleCheck,
		)

		group.POST(
			"/create_password",
			handler.middleware.CheckAccesToken(),
			handler.HandleCreatePassword,
		)

		group.GET(
			"/profile",
			handler.middleware.CheckAccesToken(),
			handler.HandleGetProfile,
		)

		group.PATCH(
			"/change_password",
		)
	}
}

func (h *ProfileHandler) GetAuthData(gctx *gin.Context) {
	authUrlResponse, err := h.ui.Usecase.Profile.CreateAuthUrlResponse()
	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, authUrlResponse)
}

func (h *ProfileHandler) HandleAuthListen(gctx *gin.Context) {
	authKey := gctx.Param("temp_id")
	if authKey == "" {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	flusher, ok := gctx.Writer.(http.Flusher)
	if !ok {
		gin_gen.HandleError(gctx, global.ErrInternalError)
		return
	}

	gctx.Writer.Header().Set("Content-Type", "text/event-stream")
	gctx.Writer.Header().Set("Cache-Control", "no-cache")
	gctx.Writer.Header().Set("Connection", "keep-alive")

	userData, err := h.ui.Usecase.Profile.WaitTgAuthVerify(gctx.Request.Context(), authKey)
	if err != nil {
		fmt.Fprint(gctx.Writer, global.SSEErrorEventMessage(global.ErrStatusCodes[err]))
		flusher.Flush()
		return
	}

	fmt.Fprint(gctx.Writer, global.SSEEventMessage(userData))
	flusher.Flush()
}

func (h *ProfileHandler) HandleTgLogin(gctx *gin.Context) {
	authID := gctx.Param("temp_id")
	deviceID := gctx.Request.Header.Get("Device-ID")

	if deviceID == "" {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	if authID == "" {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	userDataWithToken, err := h.ui.Usecase.Jwt.GenerateTokenByTempAuthData(gctx.Request.Context(), authID)
	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	err = transaction.RunInTxExec(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) error {
			return h.ui.Profile.CreateUserDeviceIDIfNotExists(ctx, userDataWithToken.UserData.ID, deviceID)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	h.setAccessToken(gctx, userDataWithToken.Token)

	gctx.JSON(http.StatusOK, userDataWithToken.UserData)
}

func (h *ProfileHandler) HandleLogin(gctx *gin.Context) {
	deviceID := gctx.Request.Header.Get("Device-ID")

	if deviceID == "" {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	var param profile.PasswordLoginParam

	if err := gctx.BindJSON(&param); err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	userData, err := transaction.RunInTx(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) (profile.PasswordLoginResponse, error) {
			return h.ui.Profile.OnPasswordLogin(ctx, param, deviceID)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	if !userData.NeedTwoStepAuth {
		token, err := h.ui.Usecase.Jwt.GenerateToken(userData.UserID, userData.Access)
		if err != nil {
			gin_gen.HandleError(gctx, err)
			return
		}
		h.setAccessToken(gctx, token)
	}

	gctx.JSON(http.StatusOK, userData)
}

func (h *ProfileHandler) HandleCheck(gctx *gin.Context) {
	token, err := gctx.Cookie("access_token")
	if err != nil {
		gin_gen.HandleError(gctx, global.ErrPermissionDenied)
		return
	}

	userData, err := h.ui.Usecase.Jwt.ParseToken(token)
	if err != nil {
		if err == global.ErrExpired {
			h.removeAccessToken(gctx, token)
		}

		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, userData)
}

func (h *ProfileHandler) HandleCreatePassword(gctx *gin.Context) {
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

	var param profile.SetPasswordParam

	if err := gctx.BindJSON(&param); err != nil {
		gin_gen.HandleError(gctx, global.ErrInvalidParam)
		return
	}

	err := transaction.RunInTxExec(
		gctx,
		h.log,
		h.sm,
		func(ctx context.Context) error {
			return h.ui.Profile.SetProfilePassword(ctx, param.Password, claimsData.UserID)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusCreated, gin.H{"success": true})
}

func (h *ProfileHandler) HandleGetProfile(gctx *gin.Context) {
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
		func(ctx context.Context) (profile.UserCommonInfo, error) {
			return h.ui.Usecase.Profile.GetUserCommonInfo(ctx, claimsData.UserID)
		},
	)

	if err != nil {
		gin_gen.HandleError(gctx, err)
		return
	}

	gctx.JSON(http.StatusCreated, data)
}

func (h *ProfileHandler) setAccessToken(gctx *gin.Context, token string) {
	lifeTime := int(h.config.JwtSecretTTL.Milliseconds())
	gctx.SetCookie("access_token", token, lifeTime, "/", "admin-panel.local", false, true)
}

func (h *ProfileHandler) removeAccessToken(gctx *gin.Context, token string) {
	gctx.SetCookie("access_token", token, -1, "/", "admin-panel.local", false, true)
}
