package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	pbe "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/delphi/aristides/proto"
	"github.com/odysseia-greek/olympia/melissos/monos"
)

const (
	DefaultPollEvery   = 10 * time.Second
	DefaultStableFor   = 1 * time.Minute
	DefaultMaxWait     = 5 * time.Minute
	DefaultMinDocs     = int64(500)
	DefaultCountReqTTL = 10 * time.Second
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

	duration := time.Millisecond * 5000
	minute := time.Minute * 60
	timeFinished := minute.Milliseconds()

	handler, err := monos.CreateNewConfig(duration, timeFinished)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	done := make(chan bool, 1) // buffered so send never blocks

	go func() {
		defer close(done)

		ctx := context.Background()
		ok := handler.WaitForDictionarySettled(
			ctx,
			DefaultMinDocs,
			DefaultPollEvery,
			DefaultStableFor,
			DefaultMaxWait,
		)

		done <- ok
	}()

	select {
	case ok := <-done:
		if ok {
			logging.Info("Dictionary settled; starting work")
		} else {
			logging.Info("Dictionary did not settle in time; aborting or fallback")
			os.Exit(1)
		}
	}

	go handler.PrintProgress()

	finishedParmenides := handler.HandleParmenides()
	if finishedParmenides {
		finishedDutch := handler.HandleDutch()
		if finishedDutch {
			logging.System("Finished Run")

			logging.Debug("setting message back so that it can be picked up by the next job")
			ctx := context.Background()
			msg := &pbe.Epistello{
				Id:      uuid.New().String(),
				Data:    "completed",
				Channel: handler.JobCompletionChannel,
			}
			_, err = handler.Eupalinos.EnqueueMessage(ctx, msg)

			logging.Debug("closing Ambassador because job is done")
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
