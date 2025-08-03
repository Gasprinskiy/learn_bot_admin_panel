package bot_api

import (
	"context"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/telegram"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/tools/bot_tool"
	"learn_bot_admin_panel/tools/logger"
	"learn_bot_admin_panel/uimport"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotProfileHandler struct {
	ui     *uimport.UsecaseImport
	b      *bot.Bot
	config *config.Config
	log    *logger.Logger
	sm     transaction.SessionManager
}

func NewBotProfileHandler(
	ui *uimport.UsecaseImport,
	b *bot.Bot,
	config *config.Config,
	log *logger.Logger,
	sm transaction.SessionManager,
) {
	handler := BotProfileHandler{
		ui,
		b,
		config,
		log,
		sm,
	}

	handler.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypePrefix,
		handler.TgAuthHandler,
	)
}

func (h *BotProfileHandler) TgAuthHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, err := transaction.RunInTxCommit(
		ctx,
		h.log,
		h.sm,
		func(ctx context.Context) (string, error) {
			return h.ui.Usecase.Profile.TgAuthVerify(ctx, update.Message.From.Username, update.Message.Text, update.Message.From.ID)
		},
	)

	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, telegram.MessagesByError[err])
		return
	}

	bot_tool.SendHTMLParseModeMessage(ctx, b, update, message)
}
