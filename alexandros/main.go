package main

import (
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/alexandros/app"
	"log"
	"net/http"
	"os"
)

const standardPort = ":5000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ALEXANDROS
	logging.System(`
  ____  _        ___  __ __   ____  ____   ___    ____   ___   _____
 /    || |      /  _]|  |  | /    ||    \ |   \  |    \ /   \ / ___/
|  o  || |     /  [_ |  |  ||  o  ||  _  ||    \ |  D  )     (   \_ 
|     || |___ |    _]|_   _||     ||  |  ||  D  ||    /|  O  |\__  |
|  _  ||     ||   [_ |     ||  _  ||  |  ||     ||    \|     |/  \ |
|  |  ||     ||     ||  |  ||  |  ||  |  ||     ||  .  \     |\    |
|__|__||_____||_____||__|__||__|__||__|__||_____||__|\_|\___/  \___|
                                                                    
`)
	logging.System("\"Ου κλέπτω την νίκην’\"")
	logging.System("\"I will not steal my victory\"")
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	env := os.Getenv("ENV")

	alexandrosConfig, err := app.CreateNewConfig(env)
	if err != nil {
		log.Print(err)
		log.Fatal("death has found me")
	}

	srv := app.InitRoutes(alexandrosConfig)

	log.Printf("%s : %s", "running on port", port)
	err = http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
