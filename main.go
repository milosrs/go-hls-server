package main

import (
	"github.com/milosrs/go-hls-server/api"
)

func main() {
	router := api.CreateRouter()
	router.Run()
}
