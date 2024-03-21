package accounting_bot

import (
	"fmt"
)

func processMessageDiagnnostics(runtime *botRuntime) string {

	var message_text string

	for k, v := range runtime.r.Header {
		message_text += fmt.Sprintf("Header: %s, value: %v\n", k, v)
	}

	message_text += fmt.Sprintf("\nYour Client IP: %s\n", runtime.r.RemoteAddr)

	return message_text

}
