package main

import (
	"os"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	Router  			*mux.Router
	ValidationRouter	*mux.Router
	bindPort			string
	//allowedOrigins		[]string
	requestTimeoutSec	int
}

func tryGetEnv(varName string, fallbackValue string) (varValue string) {
	if value, ok := os.LookupEnv(varName); ok {
		return value
	}
	return fallbackValue
}

const (
	fbURLRegex = "\\A((?:https?://)?(?:www.)?(?:facebook.com|fb.me|fb.com)(?:[^\\.])(?:/.+)?)"
 	linkedInURLRegex = "\\A((?:https?://)?(?:[www.|\\w+)])?(?:linkedin.com|lnkd.in)(?:[^\\.])(?:/.+)?)"

 	// twitterUsernameRegex = "\\A(?:@?\\w)+\\z"
	twitterURLFmtStr = "https://twitter.com/%s"

	// instagramUsernameRegex = "\\A(?:@?[\\w\\.]+)+\\z"
	instagramURLFmtStr = "https://www.instagram.com/%s/"
)

func (a *App) InitializeRouting() {
	a.Router = mux.NewRouter()
	a.Router.StrictSlash(false)
	a.ValidationRouter = a.Router.PathPrefix("/validate").Subrouter()

	log.Print("[INIT] Setting up routes")
	a.initializeApiRoutes()

	log.Print("[INIT] Initialization complete")
}

func (a *App) initializeApiRoutes() {

	a.ValidationRouter.HandleFunc("/", a.genericURLValidationView).Methods("GET")

	a.ValidationRouter.HandleFunc("/facebook", a.facebookValidationView).Methods("GET")
	a.ValidationRouter.HandleFunc("/linkedin", a.linkedInValidationView).Methods("GET")
	a.ValidationRouter.HandleFunc("/twitter", a.twitterValidationView).Methods("GET")
	a.ValidationRouter.HandleFunc("/instagram", a.instagramValidationView).Methods("GET")
}

func (a *App) Run(addr string) {

	log.Printf("Starting HTTP server at %v", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func main() {

	a := App{
		bindPort: tryGetEnv("HTTP_PORT", "9000"),
		//allowedOrigins: strings.Split(tryGetEnv("", ""), ","),
		requestTimeoutSec: 3,
	}

	a.InitializeRouting()

	a.Run(fmt.Sprintf("0.0.0.0:%s", a.bindPort))
}
