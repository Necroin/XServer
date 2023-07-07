package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"xserver/src/config"
)

const (
	errorLevel   = 0
	infoLevel    = 1
	debugLevel   = 2
	verboseLevel = 3
)

var (
	logLevel    = infoLevel
	mutex       sync.Mutex
	logLevelMap = map[string]int{
		"error":   errorLevel,
		"info":    infoLevel,
		"debug":   debugLevel,
		"verbose": verboseLevel,
	}
)

func Configure(config *config.Config) error {
	log.SetOutput(os.Stdout)
	if config.LogsPath != "" {
		if err := os.MkdirAll(path.Dir(config.LogsPath), os.ModePerm); err != nil {
			return fmt.Errorf("[XServer] [Logger] [Error] failed create logs directory: %s", err)
		}

		var logsFile *os.File
		logsFile, err := os.OpenFile(config.LogsPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			logsFile, err = os.Create(config.LogsPath)
			if err != nil {
				return fmt.Errorf("[XServer] [Logger] [Error] failed create logs file: %s", err)
			}
		}
		log.SetOutput(logsFile)
	}

	configLogLevel, ok := logLevelMap[config.LogLevel]
	if !ok {
		configLogLevel = infoLevel
	}

	logLevel = configLogLevel

	return nil
}

func Info(message string) {
	if logLevel >= infoLevel {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			log.Println("INFO: " + message)
		}()
	}
}

func Error(message string) {
	if logLevel >= errorLevel {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			log.Println("ERROR: " + message)
		}()
	}
}

func Debug(message string) {
	if logLevel >= debugLevel {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			log.Println("DEBUG: " + message)
		}()
	}
}

func Verbose(message string) {
	if logLevel >= verboseLevel {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			log.Println("VERBOSE: " + message)
		}()
	}
}
