package global

import "fmt"

const SSEMessageTemplate = "event: %s\ndata: %s\n\n"

type SSEEvent string

const (
	SSEErrorEvent SSEEvent = "error"
	SSEDoneEvent  SSEEvent = "done"
)

var (
	SSEEventMessage = func(event SSEEvent, data any) string {
		return fmt.Sprintf(SSEMessageTemplate, SSEDoneEvent, data)
	}
)
