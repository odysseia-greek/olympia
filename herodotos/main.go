package main

import (
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/herodotos/text"
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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=HERODOTOS
	logging.System(`
 __ __    ___  ____   ___   ___     ___   ______   ___   _____
|  |  |  /  _]|    \ /   \ |   \   /   \ |      | /   \ / ___/
|  |  | /  [_ |  D  )     ||    \ |     ||      ||     (   \_ 
|  _  ||    _]|    /|  O  ||  D  ||  O  ||_|  |_||  O  |\__  |
|  |  ||   [_ |    \|     ||     ||     |  |  |  |     |/  \ |
|  |  ||     ||  .  \     ||     ||     |  |  |  |     |\    |
|__|__||_____||__|\_|\___/ |_____| \___/   |__|   \___/  \___|
                                                              
`)
	logging.System("\"Ἡροδότου Ἁλικαρνησσέος ἱστορίης ἀπόδεξις ἥδε\"")
	logging.System("\"This is the display of the inquiry of Herodotos of Halikarnassos\"")
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	env := os.Getenv("ENV")

	herodotosConfig, err := text.CreateNewConfig(env)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	srv := text.InitRoutes(herodotosConfig)

	log.Printf("%s : %s", "running on port", port)
	err = http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
