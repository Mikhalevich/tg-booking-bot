package tgbot

import (
	"github.com/Mikhalevich/tg-booking-bot/internal/app/tgbot/router"
)

func Compose(routes ...RouteRegisterFunc) RouteRegisterFunc {
	return func(register router.Register) {
		for _, r := range routes {
			r(register)
		}
	}
}
