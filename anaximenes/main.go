package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"github.com/odysseia-greek/olympia/anaximenes/seeder"
	"log"
	"os"
	"strings"
)

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ANAXIMENES
	logging.System(`
  ____  ____    ____  __ __  ____  ___ ___    ___  ____     ___  _____
 /    ||    \  /    ||  |  ||    ||   |   |  /  _]|    \   /  _]/ ___/
|  o  ||  _  ||  o  ||  |  | |  | | _   _ | /  [_ |  _  | /  [_(   \_ 
|     ||  |  ||     ||_   _| |  | |  \_/  ||    _]|  |  ||    _]\__  |
|  _  ||  |  ||  _  ||     | |  | |   |   ||   [_ |  |  ||   [_ /  \ |
|  |  ||  |  ||  |  ||  |  | |  | |   |   ||     ||  |  ||     |\    |
|__|__||__|__||__|__||__|__||____||___|___||_____||__|__||_____| \___|
                                                                      
`)
	logging.System(strings.Repeat("~", 37))
	logging.System("\"οἷον ἡ ψυχή ἡ ἡμετέρα ἀὴρ οὖσα συγκρατεῖ ἡμᾶς, καὶ ὅλον τὸν κόσμον πνεῦμα καὶ ἀὴρ περιέχει\"")
	logging.System("\"Just as our soul, being air, constrains us, so breath and air envelops the whole kosmos.\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	env := os.Getenv("ENV")

	handler, err := seeder.CreateNewConfig(env)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	err = handler.DeleteIndexAtStartUp()
	if err != nil {
		logging.Debug("cannot delete index which means an aliased version exist and should not be deleted")
		logging.Error(err.Error())
		os.Exit(0)
	}

	err = handler.CreateIndexAtStartup()
	if err != nil {
		log.Fatal(err)
	}

	logging.Debug("closing ptolemaios because job is done")
	// just setting a code that could be used later to check is if it was sent from an actual service
	uuidCode := uuid.New().String()
	_, err = handler.Ambassador.ShutDown(context.Background(), &pb.ShutDownRequest{Code: uuidCode})
	if err != nil {
		logging.Error(err.Error())
	}

	os.Exit(0)

}
