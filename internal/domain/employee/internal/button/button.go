package button

import (
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func CancelButton(text string, actionData ActionData) port.Button {
	data, _ := actionData.EncodeToString()

	return port.Button{
		Text: text,
		Data: data,
	}
}
