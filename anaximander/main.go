package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"github.com/odysseia-greek/olympia/anaximander/app"
	"github.com/odysseia-greek/olympia/anaximander/config"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var documents int

//go:embed arkho
var arkho embed.FS

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ANAXIMANDER
	logging.System(`
  ____  ____    ____  __ __  ____  ___ ___   ____  ____   ___      ___  ____  
 /    ||    \  /    ||  |  ||    ||   |   | /    ||    \ |   \    /  _]|    \ 
|  o  ||  _  ||  o  ||  |  | |  | | _   _ ||  o  ||  _  ||    \  /  [_ |  D  )
|     ||  |  ||     ||_   _| |  | |  \_/  ||     ||  |  ||  D  ||    _]|    / 
|  _  ||  |  ||  _  ||     | |  | |   |   ||  _  ||  |  ||     ||   [_ |    \ 
|  |  ||  |  ||  |  ||  |  | |  | |   |   ||  |  ||  |  ||     ||     ||  .  \
|__|__||__|__||__|__||__|__||____||___|___||__|__||__|__||_____||_____||__|\_|
                                                                              
`)
	logging.System(strings.Repeat("~", 37))
	logging.System("\"οὐ γὰρ ἐν τοῖς αὐτοῖς ἐκεῖνος ἰχθῦς καὶ ἀνθρώπους, ἀλλ' ἐν ἰχθύσιν ἐγγενέσθαι τὸ πρῶτον ἀνθρώπους ἀποφαίνεται καὶ τραφέντας, ὥσπερ οἱ γαλεοί, καὶ γενομένους ἱκανους ἑαυτοῖς βοηθεῖν ἐκβῆναι τηνικαῦτα καὶ γῆς λαβέσθαι.\"")
	logging.System("\"He declares that at first human beings arose in the inside of fishes, and after having been reared like sharks, and become capable of protecting themselves, they were finally cast ashore and took to land\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	env := os.Getenv("ENV")

	anaximanderConfig, err := config.CreateNewConfig(env)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	root := "arkho"

	rootDir, err := arkho.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	handler := app.AnaximanderHandler{Config: anaximanderConfig}
	err = handler.DeleteIndexAtStartUp()
	if err != nil {
		log.Fatal(err)
	}

	err = handler.CreateIndexAtStartup()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, dir := range rootDir {
		logging.Debug("working on the following directory: " + dir.Name())
		if dir.IsDir() {
			filePath := path.Join(root, dir.Name())
			files, err := arkho.ReadDir(filePath)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				logging.Debug(fmt.Sprintf("found %s in %s", f.Name(), filePath))
				plan, _ := arkho.ReadFile(path.Join(filePath, f.Name()))
				var declension models.Declension
				err := json.Unmarshal(plan, &declension)
				if err != nil {
					log.Fatal(err)
				}

				documents += 1
				wg.Add(1)
				go func(d models.Declension) {
					defer wg.Done()
					err := handler.AddToElastic(d)
					if err != nil {
						logging.Error(err.Error())
					}
				}(declension)
			}
		}
	}

	go handler.PrintProgress(documents)
	wg.Wait()

	// Create a ticker that ticks every second
	ticker := time.NewTicker(1 * time.Second)

	// Create a channel to receive the stop signal
	stopCh := make(chan struct{})

	for i := 10; i > 0; i-- {
		logging.Debug(fmt.Sprintf("%d seconds left", i))
		// Wait for the next tick or the stop signal
		select {
		case <-ticker.C:
			// Proceed to the next second
		case <-stopCh:
			// Received the stop signal, exit the countdown
			logging.Debug("Countdown stopped.")
			ticker.Stop()
			return
		}
	}
	// Countdown reached 0, exit the program
	logging.Info(fmt.Sprintf("created: %s", strconv.Itoa(handler.Config.Created)))

	logging.Debug("closing Ptolemaios because job is done")
	// just setting a code that could be used later to check is if it was sent from an actual service
	uuidCode := uuid.New().String()
	_, err = handler.Config.Ambassador.ShutDown(context.Background(), &pb.ShutDownRequest{Code: uuidCode})
	if err != nil {
		logging.Error(err.Error())
	}
	os.Exit(0)

}
