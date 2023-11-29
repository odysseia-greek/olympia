package main

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/melissos/app"
	"github.com/odysseia-greek/olympia/melissos/config"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=MELISSOS
	logging.System(`
 ___ ___    ___  _      ____ _____ _____  ___   _____
|   |   |  /  _]| |    |    / ___// ___/ /   \ / ___/
| _   _ | /  [_ | |     |  (   \_(   \_ |     (   \_ 
|  \_/  ||    _]| |___  |  |\__  |\__  ||  O  |\__  |
|   |   ||   [_ |     | |  |/  \ |/  \ ||     |/  \ |
|   |   ||     ||     | |  |\    |\    ||     |\    |
|___|___||_____||_____||____|\___| \___| \___/  \___|
                                                     
`)
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

	duration := time.Millisecond * 5000
	minute := time.Minute * 60
	timeFinished := minute.Milliseconds()

	handler := app.MelissosHandler{
		Config:       melissosConfig,
		Duration:     duration,
		TimeFinished: timeFinished,
	}

	done := make(chan bool)

	go func() {
		handler.WaitForJobsToFinish(done)
	}()

	select {

	case <-done:
		logging.Info(fmt.Sprintf("%s job finished", melissosConfig.Job))
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
