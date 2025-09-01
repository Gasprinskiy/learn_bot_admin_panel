package chanel_kicker

type KickUserParam struct {
	TgID   int64 `form:"tg_id"`
	Reason int   `form:"reason"`
}

func NewKickUserParamWithMoneyBackReason(tgID int64, reason int) KickUserParam {
	return KickUserParam{
		TgID:   tgID,
		Reason: reason,
	}
}
