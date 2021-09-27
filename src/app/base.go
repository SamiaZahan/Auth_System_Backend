package app

import "github.com/emamulandalib/airbringr-auth/repository"

type App struct {
	Mongo *repository.Mongo
}

func New() *App {
	return &App{}
}
