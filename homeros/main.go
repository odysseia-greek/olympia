package main

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/homeros/app"
	"github.com/odysseia-greek/olympia/homeros/handlers"
	"github.com/odysseia-greek/olympia/homeros/schemas"
	"net/http"
	"os"
)

const standardPort = ":8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=HOMEROS
	logging.System(`
 __ __   ___   ___ ___    ___  ____   ___   _____
|  |  | /   \ |   |   |  /  _]|    \ /   \ / ___/
|  |  ||     || _   _ | /  [_ |  D  )     (   \_ 
|  _  ||  O  ||  \_/  ||    _]|    /|  O  |\__  |
|  |  ||     ||   |   ||   [_ |    \|     |/  \ |
|  |  ||     ||   |   ||     ||  .  \     |\    |
|__|__| \___/ |___|___||_____||__|\_|\___/  \___|
                                                 
`)
	logging.System("Αἶψα γὰρ ἐν κακότητι βροτοὶ καταγηράσκουσιν.")
	logging.System("Hardship can age a person overnight..")
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	handler := schemas.HomerosHandler()
	tracingConfig := handlers.InitTracingConfig()

	srv := app.InitRoutes(handler.Tracer, tracingConfig, handler.Randomizer)

	logging.System(fmt.Sprintf("running on port %s", port))
	err := http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
