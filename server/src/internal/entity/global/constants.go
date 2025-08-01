package global

import "fmt"

const SSEMessageTemplate = "event: %s\ndata: %s\n\n"
const SSEErrorMessageTemplate = "event: %s\ndata: %d\n\n"

type SSEEvent string

const (
	SSEErrorEvent SSEEvent = "error"
	SSEDoneEvent  SSEEvent = "user_data"
)

var (
	SSEEventMessage = func(data any) string {
		return fmt.Sprintf(SSEMessageTemplate, SSEDoneEvent, data)
	}
	SSEErrorEventMessage = func(errCode int) string {
		return fmt.Sprintf(SSEErrorMessageTemplate, SSEErrorEvent, errCode)
	}
)
