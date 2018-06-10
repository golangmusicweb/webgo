package main

import (
	//"fmt"
	"webgo/restful"
)

func main() {
	// Add a logger
	var logger restful.Logging
	logger.GetLogger()
	defer logger.Close()

	// Load config
	var config restful.Config
	config.LoadConfig()

	// get orm engine
	//dbEngine := restful.GetDbEngine("default")
}
