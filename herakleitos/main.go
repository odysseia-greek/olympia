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
	"github.com/odysseia-greek/olympia/herakleitos/seeder"
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
	logging.System(`
 __ __    ___  ____    ____  __  _  _        ___  ____  ______   ___   _____
|  |  |  /  _]|    \  /    ||  |/ ]| |      /  _]|    ||      | /   \ / ___/
|  |  | /  [_ |  D  )|  o  ||  ' / | |     /  [_  |  | |      ||     (   \_ 
|  _  ||    _]|    / |     ||    \ | |___ |    _] |  | |_|  |_||  O  |\__  |
|  |  ||   [_ |    \ |  _  ||     ||     ||   [_  |  |   |  |  |     |/  \ |
|  |  ||     ||  .  \|  |  ||  .  ||     ||     | |  |   |  |  |     |\    |
|__|__||_____||__|\_||__|__||__|\_||_____||_____||____|  |__|   \___/  \___|
                                                                            
`)
	logging.System(strings.Repeat("~", 37))
	logging.System("\"πάντα ῥεῖ\"")
	logging.System("\"everything flows\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	env := os.Getenv("ENV")

	handler, err := seeder.CreateNewConfig(env)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	root := "rhema"
	rootDir, err := rhema.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

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
	logging.Info(fmt.Sprintf("created: %s", strconv.Itoa(handler.Created)))
	logging.Info(fmt.Sprintf("texts found in rhema: %s", strconv.Itoa(documents)))

	logging.Debug("closing ptolemaios because job is done")
	// just setting a code that could be used later to check is if it was sent from an actual service
	uuidCode := uuid.New().String()
	_, err = handler.Ambassador.ShutDown(context.Background(), &pb.ShutDownRequest{Code: uuidCode})
	if err != nil {
		logging.Error(err.Error())
	}
	os.Exit(0)
}
