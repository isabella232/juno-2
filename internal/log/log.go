// Package log provides a logger.
package log

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

// Default is the default logger. It is a "sugared" variant of the zap
// logger.
var Default *zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// notest
		log.Fatalln("failed to initialise application logger")
	}
	Default = logger.Sugar()
}

// NameSlot represents a space for the Named slot in the log line.
func NameSlot(name string) string {
	return fmt.Sprintf("%-20s", name)
}
