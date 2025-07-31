package main

import (
	"github.com/yaydraco/tandem/internal/cmd"
	"github.com/yaydraco/tandem/internal/logging"
)

func main() {
	defer logging.RecoverPanic("main", func() {
		logging.ErrorPersist("Application terminated due to unhandled panic")
	})

	cmd.Execute()
}
