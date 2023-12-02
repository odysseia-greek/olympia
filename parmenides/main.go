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
	"github.com/odysseia-greek/olympia/parmenides/app"
	"github.com/odysseia-greek/olympia/parmenides/config"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

//go:embed sullego
var sullego embed.FS

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PARMENIDES
	logging.System("\n ____   ____  ____   ___ ___    ___  ____   ____  ___      ___  _____\n|    \\ /    ||    \\ |   |   |  /  _]|    \\ |    ||   \\    /  _]/ ___/\n|  o  )  o  ||  D  )| _   _ | /  [_ |  _  | |  | |    \\  /  [_(   \\_ \n|   _/|     ||    / |  \\_/  ||    _]|  |  | |  | |  D  ||    _]\\__  |\n|  |  |  _  ||    \\ |   |   ||   [_ |  |  | |  | |     ||   [_ /  \\ |\n|  |  |  |  ||  .  \\|   |   ||     ||  |  | |  | |     ||     |\\    |\n|__|  |__|__||__|\\_||___|___||_____||__|__||____||_____||_____| \\___|\n                                                                     \n")
	logging.System(strings.Repeat("~", 37))
	logging.System("\"τό γάρ αυτο νοειν έστιν τε καί ειναι\"")
	logging.System("\"for it is the same thinking and being\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	env := os.Getenv("ENV")

	parmenidesConfig, conn, err := config.CreateNewConfig(env)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	defer conn.Close()

	root := "sullego"
	rootDir, err := sullego.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	handler := app.ParmenidesHandler{Config: parmenidesConfig}

	err = handler.DeleteIndexAtStartUp()
	if err != nil {
		log.Fatal(err)
	}
	err = handler.CreateIndexAtStartup()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	documents := 0

	for _, dir := range rootDir {
		logging.Debug("working on the following directory: " + dir.Name())
		if dir.IsDir() {
			method := dir.Name()
			logging.Info(fmt.Sprintf("working on %s", method))
			methodPath := path.Join(root, dir.Name())
			methodDir, err := sullego.ReadDir(methodPath)
			if err != nil {
				log.Fatal(err)
			}

			for _, innerDir := range methodDir {
				category := innerDir.Name()
				filePath := path.Join(root, dir.Name(), innerDir.Name())
				files, err := sullego.ReadDir(filePath)
				if err != nil {
					log.Fatal(err)
				}
				for _, f := range files {
					logging.Debug(fmt.Sprintf("found %s in %s", f.Name(), filePath))
					plan, _ := sullego.ReadFile(path.Join(filePath, f.Name()))
					var logoi models.Logos
					err := json.Unmarshal(plan, &logoi)
					if err != nil {
						log.Fatal(err)
					}

					logging.Info(fmt.Sprintf("method: %s | category: %s | documents: %d", method, category, len(logoi.Logos)))

					documents += len(logoi.Logos)

					wg.Add(1)
					go func() {
						err := handler.Add(logoi, &wg, method, category)
						if err != nil {
							log.Fatal(err)
						}
					}()
				}
			}
		}

	}

	wg.Wait()
	logging.Info(fmt.Sprintf("created: %s", strconv.Itoa(parmenidesConfig.Created)))
	logging.Info(fmt.Sprintf("words found in sullego: %s", strconv.Itoa(documents)))

	logging.Debug("closing Ptolemaios because job is done")
	// just setting a code that could be used later to check is if it was sent from an actual service
	uuidCode := uuid.New().String()
	_, err = handler.Config.Ambassador.ShutDown(context.Background(), &pb.ShutDownRequest{Code: uuidCode})
	if err != nil {
		logging.Error(err.Error())
	}

	os.Exit(0)
}
