package main

import (
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/melissos/app"
	"github.com/odysseia-greek/olympia/melissos/config"
	"log"
	"os"
	"strings"
)

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=MELISSOS
	logging.System("\n ___ ___    ___  _      ____ _____ _____  ___   _____\n|   |   |  /  _]| |    |    / ___// ___/ /   \\ / ___/\n| _   _ | /  [_ | |     |  (   \\_(   \\_ |     (   \\_ \n|  \\_/  ||    _]| |___  |  |\\__  |\\__  ||  O  |\\__  |\n|   |   ||   [_ |     | |  |/  \\ |/  \\ ||     |/  \\ |\n|   |   ||     ||     | |  |\\    |\\    ||     |\\    |\n|___|___||_____||_____||____|\\___| \\___| \\___/  \\___|\n                                                     \n")
	logging.System(strings.Repeat("~", 37))
	logging.System("\"Οὕτως οὖν ἀίδιόν ἐστι καὶ ἄπειρον καὶ ἓν καὶ ὅμοιον πᾶν.\"")
	logging.System("\"So then it is eternal and infinite and one and all alike.\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	env := os.Getenv("ENV")

	melissosConfig, conn, err := config.CreateNewConfig(env)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	handler := app.MelissosHandler{
		Config: melissosConfig,
	}

	go handler.PrintProgress()

	finishedParmenides := handler.HandleParmenides()
	if finishedParmenides {
		finishedDutch := handler.HandleDutch()
		if finishedDutch {
			logging.System("Finished Run")
			conn.Close()
			os.Exit(0)
		}
	}
}
