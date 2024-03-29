package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"github.com/odysseia-greek/olympia/melissos/ergon"
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
	duration := time.Millisecond * 5000
	minute := time.Minute * 60
	timeFinished := minute.Milliseconds()

	handler, conn, err := ergon.CreateNewConfig(env, duration, timeFinished)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	done := make(chan bool)

	go func() {
		handler.WaitForJobsToFinish(done)
	}()

	select {

	case <-done:
		logging.Info(fmt.Sprintf("%s job finished", handler.Job))
	}

	go handler.PrintProgress()

	finishedParmenides := handler.HandleParmenides()
	if finishedParmenides {
		finishedDutch := handler.HandleDutch()
		if finishedDutch {
			logging.System("Finished Run")
			conn.Close()

			logging.Debug("closing Ptolemaios because job is done")
			// just setting a code that could be used later to check is if it was sent from an actual service
			uuidCode := uuid.New().String()
			_, err = handler.Ambassador.ShutDown(context.Background(), &pb.ShutDownRequest{Code: uuidCode})
			if err != nil {
				logging.Error(err.Error())
			}
			os.Exit(0)
		}
	}
}
