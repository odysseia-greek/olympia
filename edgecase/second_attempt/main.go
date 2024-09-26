package main

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/anotherone/api"
	"log"
	"net/http"
	"os"
)

const standardPort = ":5000"

func main() {
	port := os.Getenv("5000")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=diogenes
	logging.System(`
 ___    ____  ___    ____    ___  ____     ___  _____
|   \  |    |/   \  /    |  /  _]|    \   /  _]/ ___/
|    \  |  ||     ||   __| /  [_ |  _  | /  [_(   \_
|  D  | |  ||  O  ||  |  ||    _]|  |  ||    _]\__  |
|     | |  ||     ||  |_ ||   [_ |  |  ||   [_ /  \ |
|     | |  ||     ||     ||     ||  |  ||     |\    |
|_____||____|\___/ |___,_||_____||__|__||_____| \___|
`)
	logging.System("\"Ἀποσκότησόν μοι\"")
	logging.System("\"Get out of my light.\"")
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	apiConfig, err := api.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	srv := api.InitRoutes(apiConfig)

	logging.Info(fmt.Sprintf("%s : %s", "running on port", port))
	err = http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
