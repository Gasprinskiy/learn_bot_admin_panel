package bot_tool

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SendHTMLParseModeMessage(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	message string,
) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      message,
		ParseMode: "HTML",
	})
}

// func SendReplyKeyboardMessage(
// 	ctx context.Context,
// 	b *bot.Bot,
// 	update *models.Update,
// 	replyMessage global.ReplyMessage,
// 	closeAfterClick bool,
// ) {
// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: update.Message.Chat.ID,
// 		Text:   replyMessage.Message,
// 		ReplyMarkup: &models.ReplyKeyboardMarkup{
// 			Keyboard: [][]models.KeyboardButton{
// 				replyMessage.ButtonList,
// 			},
// 			ResizeKeyboard:  true,
// 			OneTimeKeyboard: closeAfterClick,
// 		},
// 	})
// }
