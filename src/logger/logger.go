package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"xserver/src/config"
)

const (
	errorLevel   = 0
	infoLevel    = 1
	debufLevel   = 2
	verboseLevel = 3
)

var (
	logLevel = infoLevel
)

func Configure(config *config.Config) error {
	log.SetOutput(os.Stdout)
	if config.LogsPath != "" {
		if err := os.MkdirAll(path.Dir(config.LogsPath), os.ModePerm); err != nil {
			return fmt.Errorf("[XServer] [Logger] [Error] failed create logs directory: %s", err)
		}

		logsFile, err := os.Create(config.LogsPath)
		if err != nil {
			return fmt.Errorf("[XServer] [Logger] [Error] failed create logs file: %s", err)
		}
		log.SetOutput(logsFile)
	}
	return nil
}

func Info(message string) {
	if logLevel >= infoLevel {
		go func() { log.Println("INFO: " + message) }()
	}
}

func Error(message string) {
	if logLevel >= errorLevel {
		go func() { log.Println("ERROR: " + message) }()
	}
}

func Debug(message string) {
	if logLevel >= debufLevel {
		go func() { log.Println("DEBUG: " + message) }()
	}
}

func Verbose(message string) {
	if logLevel >= verboseLevel {
		go func() { log.Println("VERBOSE: " + message) }()
	}
}
