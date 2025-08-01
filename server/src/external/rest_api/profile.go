package rest_api

import (
	"fmt"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/tools/gin_gen"

	"learn_bot_admin_panel/uimport"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	ui     *uimport.UsecaseImport
	router *gin.RouterGroup
	config *config.Config
}

func NewProfileHandler(
	ui *uimport.UsecaseImport,
	router *gin.RouterGroup,
	config *config.Config,
) {
	handler := ProfileHandler{
		ui,
		router,
		config,
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

		group.POST(
			"/login",
			handler.HandleLogin,
		)

		group.GET(
			"/check",
			handler.HandleCheck,
		)

		group.POST(
			"/create_password",
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

func (h *ProfileHandler) HandleLogin(gctx *gin.Context) {
	authID := gctx.Query("temp_id")
	if authID != "" {
		userDataWithToken, err := h.ui.Usecase.Jwt.GenerateTokenByTempAuthData(gctx.Request.Context(), authID)
		if err != nil {
			gin_gen.HandleError(gctx, err)
			return
		}

		h.setAccessToken(gctx, userDataWithToken.Token)

		gctx.JSON(http.StatusOK, userDataWithToken.UserData)
	}
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

func (h *ProfileHandler) setAccessToken(gctx *gin.Context, token string) {
	lifeTime := int(h.config.JwtSecretTTL.Milliseconds())
	gctx.SetCookie("access_token", token, lifeTime, "/", "admin-panel.local", false, true)
}

func (h *ProfileHandler) removeAccessToken(gctx *gin.Context, token string) {
	gctx.SetCookie("access_token", token, -1, "/", "admin-panel.local", false, true)
}
