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
