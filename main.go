package main

import (
	"os"
	"fmt"
	"log"
	"github.com/gorilla/mux"
)

type App struct {
	Router  			*mux.Router
	ValidationRouter	*mux.Router
	bindPort			string
}

func tryGetEnv(varName string, fallbackValue string) (varValue string) {
	if value, ok := os.LookupEnv(varName); ok {
		return value
	}
	return fallbackValue
}

func (a *App) InitializeRouting() {
	a.Router = mux.NewRouter()
	a.Router.StrictSlash(true)
	a.ValidationRouter = a.Router.PathPrefix("/val").Subrouter()

	log.Print("[INIT] Setting up routes")
	a.initializeApiRoutes()

	log.Print("[INIT] Initialization complete")
}

func (a *App) initializeApiRoutes() {

	a.ValidationRouter.HandleFunc("/jobs", a.getJobsRoute).Methods("GET")
}

func main() {

	a := App{
		bindPort: tryGetEnv("HTTP_PORT", "9000"),
	}

	a.InitializeRouting()

	a.Run(fmt.Sprintf("0.0.0.0:%s", a.bindPort))
}
