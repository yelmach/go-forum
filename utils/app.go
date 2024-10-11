package utils

type App struct {
	routers []*Router
}

func NewApp() *App {
	return &App{}
}
