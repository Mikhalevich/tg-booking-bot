package tgbot

func Compose(routes ...RouteRegisterFunc) RouteRegisterFunc {
	return func(register Register) {
		for _, r := range routes {
			r(register)
		}
	}
}
