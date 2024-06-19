package main

import (
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/protagoras/seeder"
	"log"
	"os"
)

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=protagoras
	logging.System(`
 ____  ____   ___   ______   ____   ____   ___   ____    ____  _____
|    \|    \ /   \ |      | /    | /    | /   \ |    \  /    |/ ___/
|  o  )  D  )     ||      ||  o  ||   __||     ||  D  )|  o  (   \_ 
|   _/|    /|  O  ||_|  |_||     ||  |  ||  O  ||    / |     |\__  |
|  |  |    \|     |  |  |  |  _  ||  |_ ||     ||    \ |  _  |/  \ |
|  |  |  .  \     |  |  |  |  |  ||     ||     ||  .  \|  |  |\    |
|__|  |__|\_|\___/   |__|  |__|__||___,_| \___/ |__|\_||__|__| \___|
                                                                    
`)
	logging.System("\"Πάντων χρημάτων μέτρον ἐστὶν ἄνθρωπος.\"")
	logging.System("\"Man is the measure of all things.\"")
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	handler, err := seeder.CreateNewConfig()
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	err = handler.Start()
	if err != nil {
		logging.Error(err.Error())
	}

	os.Exit(0)
}
