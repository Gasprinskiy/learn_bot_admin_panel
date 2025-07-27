package telegram

import "learn_bot_admin_panel/internal/entity/global"

var MessagesByError = map[error]string{
	global.ErrInternalError:    "Внутреняя ошибка бота, попробуйте позже или свяжитесь с поддержкой.",
	global.ErrPermissionDenied: "Отказано в доступе выполнения команды",
	global.ErrExpired:          "Срок сессии истек, попробуйте еще раз",
	global.ErrNoData:           "Данные пользователя не найдены",
}
