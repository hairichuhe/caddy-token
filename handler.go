package token

type Handler struct {
	Next   httpserver.Handler
	Config HandlerConfiguration
}
