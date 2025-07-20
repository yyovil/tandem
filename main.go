package main

import (
	"github.com/yyovil/tandem/internal/cmd"
	"github.com/yyovil/tandem/internal/logging"
)

func main() {
	defer logging.RecoverPanic("main", func() {
		logging.ErrorPersist("Application terminated due to unhandled panic")
	})

	cmd.Execute()
}
