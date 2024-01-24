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
	"github.com/odysseia-greek/olympia/parmenides/seeder"
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
	logging.System(`
 ____   ____  ____   ___ ___    ___  ____   ____  ___      ___  _____
|    \ /    ||    \ |   |   |  /  _]|    \ |    ||   \    /  _]/ ___/
|  o  )  o  ||  D  )| _   _ | /  [_ |  _  | |  | |    \  /  [_(   \_ 
|   _/|     ||    / |  \_/  ||    _]|  |  | |  | |  D  ||    _]\__  |
|  |  |  _  ||    \ |   |   ||   [_ |  |  | |  | |     ||   [_ /  \ |
|  |  |  |  ||  .  \|   |   ||     ||  |  | |  | |     ||     |\    |
|__|  |__|__||__|\_||___|___||_____||__|__||____||_____||_____| \___|
                                                                     
`)
	logging.System(strings.Repeat("~", 37))
	logging.System("\"τό γάρ αυτο νοειν έστιν τε καί ειναι\"")
	logging.System("\"for it is the same thinking and being\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	env := os.Getenv("ENV")

	handler, conn, err := seeder.CreateNewConfig(env)
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
		if dir.IsDir() {
			typePath := path.Join(root, dir.Name())
			typeDir, err := sullego.ReadDir(typePath)
			if err != nil {
				log.Fatal(err)
			}

			for _, quizType := range typeDir {
				quizPath := path.Join(typePath, quizType.Name())
				content, err := sullego.ReadFile(quizPath)
				if err != nil {
					logging.Error(err.Error())
					continue
				}

				logging.Debug(fmt.Sprintf("working on file: %s in quiz: %s", quizPath, dir.Name()))

				switch dir.Name() {
				case "media":
					wg.Add(1)
					go func() {
						err := handler.AddWithoutQueue(content, &wg)
						if err != nil {
							logging.Error(err.Error())
						}
					}()
				case "dialogue":
					wg.Add(1)
					go func() {
						var quiz []models.DialogueQuiz
						err = json.Unmarshal(content, &quiz)
						if err != nil {
							logging.Error(err.Error())
						}

						for _, q := range quiz {
							asJson, err := json.Marshal(q)
							if err != nil {
								logging.Error(err.Error())
								continue
							}

							err = handler.AddWithoutQueue(asJson, &wg)
							if err != nil {
								logging.Error(err.Error())
							}
						}
					}()
				case "authorbased":
					wg.Add(1)
					var quiz []models.AuthorBasedQuiz
					err := json.Unmarshal(content, &quiz)
					if err != nil {
						logging.Error(err.Error())
						return
					}

					err = handler.AddWithQueue(quiz, &wg)
					if err != nil {
						logging.Error(err.Error())
					}
				}
			}

		}
	}

	wg.Wait()
	logging.Info(fmt.Sprintf("created: %s", strconv.Itoa(handler.Created)))
	logging.Info(fmt.Sprintf("words found in sullego: %s", strconv.Itoa(documents)))

	logging.Debug("closing ptolemaios because job is done")
	// just setting a code that could be used later to check is if it was sent from an actual service
	uuidCode := uuid.New().String()
	_, err = handler.Ambassador.ShutDown(context.Background(), &pb.ShutDownRequest{Code: uuidCode})
	if err != nil {
		logging.Error(err.Error())
	}

	os.Exit(0)
}
