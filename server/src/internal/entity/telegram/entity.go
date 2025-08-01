package telegram

import "fmt"

type BotInfoResponse struct {
	Ok     bool       `json:"ok"`
	Result BotDetails `json:"result"`
}

type BotDetails struct {
	ID                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	Username                string `json:"username"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
	CanConnectToBusiness    bool   `json:"can_connect_to_business"`
	HasMainWebApp           bool   `json:"has_main_web_app"`
}

func (b BotDetails) BotStartUrlWithQuery(query string) string {
	return fmt.Sprintf("%s?start=%s", AccountUrl(b.Username), query)
}
