package main

import "github.com/sbxb/loyalty/internal/logger"

func main() {
	logger.SetLevel("DEBUG")

	logger.Info("Hello world")
}
