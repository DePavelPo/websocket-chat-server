package handler

type Handler struct {
	allowedOrigins []string
}

func NewHandler(origins []string) *Handler {
	return &Handler{
		allowedOrigins: origins,
	}
}
