package http_rep

import (
	"fmt"
	"learn_bot_admin_panel/internal/entity/telegram"
	"learn_bot_admin_panel/internal/repository"
	"learn_bot_admin_panel/tools/http_req"
)

type tgBot struct {
	apiPath  string
	botToken string
}

func NewTgBot(apiPath string, botToken string) repository.TgBot {
	return &tgBot{apiPath, botToken}
}

func (r *tgBot) GetBotInfo() (telegram.BotInfoResponse, error) {
	url := fmt.Sprintf("%s/bot%s/getMe", r.apiPath, r.botToken)

	return http_req.Get[telegram.BotInfoResponse](url)
}
