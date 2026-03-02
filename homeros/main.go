package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/homeros/gateway"
	"github.com/odysseia-greek/olympia/homeros/routing"
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

	tracingConfig := gateway.InitTracingConfig()
	handler, err := gateway.CreateNewConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	graphqlServer := routing.InitRoutes(handler, tracingConfig, handler.Randomizer)

	logging.System(fmt.Sprintf("running on port %s", port))
	err = http.ListenAndServe(port, graphqlServer)
	if err != nil {
		panic(err)
	}
}
