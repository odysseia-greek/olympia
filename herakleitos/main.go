package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/olympia/herakleitos/app"
	"github.com/odysseia-greek/olympia/herakleitos/config"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

//go:embed rhema
var rhema embed.FS

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=HERAKLEITOS
	logging.System("\n __ __    ___  ____    ____  __  _  _        ___  ____  ______   ___   _____\n|  |  |  /  _]|    \\  /    ||  |/ ]| |      /  _]|    ||      | /   \\ / ___/\n|  |  | /  [_ |  D  )|  o  ||  ' / | |     /  [_  |  | |      ||     (   \\_ \n|  _  ||    _]|    / |     ||    \\ | |___ |    _] |  | |_|  |_||  O  |\\__  |\n|  |  ||   [_ |    \\ |  _  ||     ||     ||   [_  |  |   |  |  |     |/  \\ |\n|  |  ||     ||  .  \\|  |  ||  .  ||     ||     | |  |   |  |  |     |\\    |\n|__|__||_____||__|\\_||__|__||__|\\_||_____||_____||____|  |__|   \\___/  \\___|\n                                                                            \n")
	logging.System(strings.Repeat("~", 37))
	logging.System("\"πάντα ῥεῖ\"")
	logging.System("\"everything flows\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	env := os.Getenv("ENV")

	herakleitosConfig, err := config.CreateNewConfig(env)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	root := "rhema"
	rootDir, err := rhema.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	handler := app.HerakleitosHandler{Config: herakleitosConfig}

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
			filePath := path.Join(root, dir.Name())
			files, err := rhema.ReadDir(filePath)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				logging.Debug(fmt.Sprintf("found %s in %s", f.Name(), filePath))
				plan, _ := rhema.ReadFile(path.Join(filePath, f.Name()))
				var rhemai models.Rhema
				err := json.Unmarshal(plan, &rhemai)
				if err != nil {
					log.Fatal(err)
				}

				documents += len(rhemai.Rhemai)

				wg.Add(1)
				go func() {
					err := handler.Add(rhemai, &wg)
					if err != nil {
						logging.Error(err.Error())
					}
				}()
			}
		}
	}

	wg.Wait()
	logging.Info(fmt.Sprintf("created: %s", strconv.Itoa(handler.Config.Created)))
	logging.Info(fmt.Sprintf("texts found in rhema: %s", strconv.Itoa(documents)))
	os.Exit(0)
}
