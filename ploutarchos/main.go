package main

import (
	"github.com/odysseia-greek/olympia/ploutarchos/app"
	"log"
	"net/http"
	"os"
)

const standardPort = ":3000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PLOUTARCHOS
	log.Print(`
 ____  _       ___   __ __  ______   ____  ____      __  __ __   ___   _____
|    \| |     /   \ |  |  ||      | /    ||    \    /  ]|  |  | /   \ / ___/
|  o  ) |    |     ||  |  ||      ||  o  ||  D  )  /  / |  |  ||     (   \_ 
|   _/| |___ |  O  ||  |  ||_|  |_||     ||    /  /  /  |  _  ||  O  |\__  |
|  |  |     ||     ||  :  |  |  |  |  _  ||    \ /   \_ |  |  ||     |/  \ |
|  |  |     ||     ||     |  |  |  |  |  ||  .  \\     ||  |  ||     |\    |
|__|  |_____| \___/  \__,_|  |__|  |__|__||__|\_| \____||__|__| \___/  \___|
                                                                            
`)
	log.Print("starting up.....")
	log.Print("starting up and getting env variables")

	srv := app.InitRoutes()

	log.Printf("%s : %s", "running on port", port)
	err := http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
