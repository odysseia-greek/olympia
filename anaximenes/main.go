package main

import (
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/anaximenes/app"
	"github.com/odysseia-greek/olympia/anaximenes/config"
	"log"
	"os"
	"strings"
)

var documents int

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ANAXIMENES
	logging.System(`
  ____  ____    ____  __ __  ____  ___ ___    ___  ____     ___  _____
 /    ||    \  /    ||  |  ||    ||   |   |  /  _]|    \   /  _]/ ___/
|  o  ||  _  ||  o  ||  |  | |  | | _   _ | /  [_ |  _  | /  [_(   \_ 
|     ||  |  ||     ||_   _| |  | |  \_/  ||    _]|  |  ||    _]\__  |
|  _  ||  |  ||  _  ||     | |  | |   |   ||   [_ |  |  ||   [_ /  \ |
|  |  ||  |  ||  |  ||  |  | |  | |   |   ||     ||  |  ||     |\    |
|__|__||__|__||__|__||__|__||____||___|___||_____||__|__||_____| \___|
                                                                      
`)
	logging.System(strings.Repeat("~", 37))
	logging.System("\"οἷον ἡ ψυχή ἡ ἡμετέρα ἀὴρ οὖσα συγκρατεῖ ἡμᾶς, καὶ ὅλον τὸν κόσμον πνεῦμα καὶ ἀὴρ περιέχει\"")
	logging.System("\"Just as our soul, being air, constrains us, so breath and air envelops the whole kosmos.\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	env := os.Getenv("ENV")

	anaximenesConfig, err := config.CreateNewConfig(env)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	handler := app.AnaximenesConfig{Config: anaximenesConfig}
	err = handler.DeleteIndexAtStartUp()
	if err != nil {
		logging.Debug("cannot delete index which means an aliased version exist and should not be deleted")
		logging.Error(err.Error())
		os.Exit(0)
	}

	err = handler.CreateIndexAtStartup()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)

}
