package telegram

import (
	"fmt"
)

const (
	DefaultURL = "https://t.me"
)

var (
	AccountUrl = func(userName string) string {
		return fmt.Sprintf("%s/%s", DefaultURL, userName)
	}
)

const (
	TwoStepAuthCallBackQueryAccept  = "tow_step_auth:accept"
	TwoStepAuthCallBackQueryDecline = "tow_step_auth:yes"
)

var TwoStepAuthCallBackQueryButtonsMap = map[string]string{
	TwoStepAuthCallBackQueryAccept:  "✅ Да",
	TwoStepAuthCallBackQueryDecline: "❌ Нет",
}

const (
	AuthSuccessfulyMessage = "Здравствуйте %s, Авторизация прошла успешно, можете вернутся обратно!"
	TwoStepAuthMessage     = "<b>%s.</b>,<br> Обнаружен вход с нового устройства, это были вы?"
)
