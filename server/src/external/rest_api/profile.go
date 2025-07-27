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
	gin    *gin.Engine
	config *config.Config
}

func NewProfileHandler(
	ui *uimport.UsecaseImport,
	gin *gin.Engine,
	config *config.Config,
) {
	handler := ProfileHandler{
		ui,
		gin,
		config,
	}

	group := handler.gin.Group("/auth")

	{
		group.GET(
			"/temp_data",
			handler.GetAuthData,
		)

		group.GET(
			"/listen",
			handler.HandleAuthListen,
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
	authKey := gctx.Query("auth_id")
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
		fmt.Fprintf(gctx.Writer, "event: error\ndata: %s\n\n", err.Error())
		flusher.Flush()
		return
	}

	fmt.Fprintf(gctx.Writer, "event: done\ndata: %s\n\n", userData)
	flusher.Flush()
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
