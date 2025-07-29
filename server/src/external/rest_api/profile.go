package rest_api

import (
	"fmt"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/global"

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
		gctx.JSON(global.ErrStatusCodes[err], gin.H{"message": err.Error()})
		return
	}

	gctx.JSON(http.StatusOK, authUrlResponse)
}

func (h *ProfileHandler) HandleAuthListen(gctx *gin.Context) {
	authKey := gctx.Param("temp_id")
	if authKey == "" {
		gctx.JSON(http.StatusBadRequest, gin.H{"message": global.ErrInvalidParam})
		return
	}

	flusher, ok := gctx.Writer.(http.Flusher)
	if !ok {
		gctx.JSON(http.StatusInternalServerError, gin.H{"message": global.ErrInternalError})
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
			gctx.JSON(global.ErrStatusCodes[err], gin.H{"message": err.Error()})
		}

		h.setAccessToken(gctx, userDataWithToken.Token)

		gctx.JSON(http.StatusOK, userDataWithToken.UserData)
	}
}

func (h *ProfileHandler) setAccessToken(gctx *gin.Context, token string) {
	lifeTime := int(h.config.JwtSecretTTL.Milliseconds())
	gctx.SetCookie("access_token", token, lifeTime, "/", "", true, true)
}

// authSession, exists := h.authChan.Read(authID)
// if !exists {
// 	gctx.JSON(http.StatusGone, gin.H{"message": global.ErrExpired})
// 	return
// }

// defer h.authChan.CleanUp(authID)

// select {
// case <-gctx.Request.Context().Done():
// 	return

// case userData := <-authSession.Chan:
// 	jsonData, err := json.Marshal(userData)
// 	if err != nil {
// 		gctx.JSON(http.StatusInternalServerError, gin.H{"message": global.ErrInternalError})
// 		return
// 	}

// 	fmt.Fprintf(gctx.Writer, "event: done\ndata: %s\n\n", jsonData)
// 	flusher.Flush()
// 	return

// case <-time.After(h.config.SSETTL):
// 	gctx.JSON(http.StatusGone, gin.H{"message": global.ErrExpired})
// 	return
// }
