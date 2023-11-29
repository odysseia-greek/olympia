package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/olympia/demokritos/app"
	"github.com/odysseia-greek/olympia/demokritos/config"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

var documents int

//go:embed lexiko
var lexiko embed.FS

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=DEMOKRITOS
	logging.System(`
 ___      ___  ___ ___   ___   __  _  ____   ____  ______   ___   _____
|   \    /  _]|   |   | /   \ |  |/ ]|    \ |    ||      | /   \ / ___/
|    \  /  [_ | _   _ ||     ||  ' / |  D  ) |  | |      ||     (   \_ 
|  D  ||    _]|  \_/  ||  O  ||    \ |    /  |  | |_|  |_||  O  |\__  |
|     ||   [_ |   |   ||     ||     ||    \  |  |   |  |  |     |/  \ |
|     ||     ||   |   ||     ||  .  ||  .  \ |  |   |  |  |     |\    |
|_____||_____||___|___| \___/ |__|\_||__|\_||____|  |__|   \___/  \___|
                                                                       
`)
	logging.System(strings.Repeat("~", 37))
	logging.System("\"νόμωι (γάρ φησι) γλυκὺ καὶ νόμωι πικρόν, νόμωι θερμόν, νόμωι ψυχρόν, νόμωι χροιή, ἐτεῆι δὲ ἄτομα καὶ κενόν\"")
	logging.System("\"By convention sweet is sweet, bitter is bitter, hot is hot, cold is cold, color is color; but in truth there are only atoms and the void.\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	env := os.Getenv("ENV")

	demokritosConfig, err := config.CreateNewConfig(env)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	root := "lexiko"

	rootDir, err := lexiko.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	handler := app.DemokritosHandler{Config: demokritosConfig}

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
			files, err := lexiko.ReadDir(filePath)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				logging.Debug(fmt.Sprintf("found %s in %s", f.Name(), filePath))
				plan, _ := lexiko.ReadFile(path.Join(filePath, f.Name()))
				var biblos models.Biblos
				err := json.Unmarshal(plan, &biblos)
				if err != nil {
					log.Fatal(err)
				}

				documents += len(biblos.Biblos)

				wg.Add(1)
				go handler.AddDirectoryToElastic(biblos, &wg)
			}
		}
	}

	wg.Wait()

	logging.Info(fmt.Sprintf("created: %s", strconv.Itoa(handler.Config.Created)))
	logging.Info(fmt.Sprintf("words found in sullego: %s", strconv.Itoa(documents)))
	os.Exit(0)
}
