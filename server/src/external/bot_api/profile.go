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

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		telegram.TwoStepAuthPrefix,
		bot.MatchTypePrefix,
		handler.TwoStepAuthHandler,
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

func (h *BotProfileHandler) TwoStepAuthHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, err := transaction.RunInTxCommit(
		ctx,
		h.log,
		h.sm,
		func(ctx context.Context) (string, error) {
			return h.ui.Usecase.Profile.TwoStepTgAuthVerify(ctx, update.CallbackQuery.From.Username, update.CallbackQuery.Data, update.CallbackQuery.From.ID)
		},
	)

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            "✅",
	})

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.From.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})

	if err != nil {
		bot_tool.SendHTMLParseModeMessageFromCallBackQuery(ctx, b, update, telegram.MessagesByError[err])
		return
	}

	bot_tool.SendHTMLParseModeMessageFromCallBackQuery(ctx, b, update, message)
}
