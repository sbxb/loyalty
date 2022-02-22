package main

import (
	"github.com/sbxb/loyalty/api"
	"github.com/sbxb/loyalty/internal/logger"
)

func main() {
	logger.SetLevel("DEBUG")

	router := api.NewRouter()
	_ = router
}
