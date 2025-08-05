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

const TwoStepAuthPrefix = "tow_step_auth"

const (
	TwoStepAuthAccept  = "accept"
	TwoStepAuthDecline = "decline"
)

var (
	TwoStepAuthCallBackQueryAccept = func(uuID string) string {
		return fmt.Sprintf("%s:%s:%s", TwoStepAuthPrefix, TwoStepAuthAccept, uuID)
	}
	TwoStepAuthCallBackQueryDecline = func(uuID string) string {
		return fmt.Sprintf("%s:%s:%s", TwoStepAuthPrefix, TwoStepAuthDecline, uuID)
	}
)

const (
	TwoStepAuthCallBackQueryAcceptView  = "✅ Да"
	TwoStepAuthCallBackQueryDeclineView = "❌ Нет"
)

const (
	AuthSuccessfulyMessage    = "Здравствуйте %s! Авторизация прошла успешно, можете вернутся обратно!"
	TwoStepAuthMessage        = "%s! Обнаружен вход с нового устройства, это были вы?"
	TwoStepAuthDeclineMessage = "⚠️ Спасибо за подтверждение.\n\nМы зафиксировали подозрительную попытку входа в ваш аккаунт и заблокировали её.\n\nВ целях безопасности мы рекомендуем cменить пароль от аккаунта"
)
