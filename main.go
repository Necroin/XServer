package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"xserver/src/builders"
	"xserver/src/config"
	"xserver/src/logger"
	"xserver/src/runners"
	"xserver/src/server"
)

var (
	commands = map[string]func(config *config.Config) error{
		"create": create,
		"build":  build,
		"start":  start,
	}
	configPath             = "./config.yml"
	handlersFilesPath      = "bin/handlers/"
	languagesBuildCommands = map[string]func(string, string) error{
		".go": builders.Go,
	}
	languagesRunCommands = map[string]func(string, http.ResponseWriter, *http.Request, func(string, error), func(string)){
		".go": runners.Executable,
	}
)

func create(config *config.Config) error {
	logger.Info("XServer: create project.")
	return nil
}

func build(config *config.Config) error {
	logger.Info("[XServer] [Build] Build project")
	if err := os.MkdirAll(handlersFilesPath, os.ModePerm); err != nil {
		return fmt.Errorf("[XServer] [Build] [Error] failed create files directory: %s", err)
	}

	for handlerName, handler := range config.Handlers {
		logger.Info(fmt.Sprintf(`[XServer] [Build] build "%s" handler`, handlerName))
		buildCommand, ok := languagesBuildCommands[path.Ext(handler.File)]
		if ok {
			if err := buildCommand(handler.File, path.Join(handlersFilesPath, handlerName, "executable")); err != nil {
				return fmt.Errorf(`[XServer] [Build] [Error] failed compile "%s" handler: %s`, handlerName, err)
			}
		}
	}
	return nil
}

func start(config *config.Config) error {
	logger.Info("[XServer] Start project")
	xserver := server.Create(
		&server.ServerOptions{Url: config.Url},
	)

	for handlerName, handler := range config.Handlers {
		xserver.AddHandler(
			handler.Path,
			func(writer http.ResponseWriter, request *http.Request) {
				logger.Info(fmt.Sprintf("[XServer] [%s Handler] handler called", handlerName))
				runCommand, ok := languagesRunCommands[path.Ext(handler.File)]
				if !ok {
					message := fmt.Sprintf("[XServer] [%s Handler] [Error] run command is unknown", handlerName)
					logger.Error(message)
					writer.Write([]byte(message))
					return
				}

				_, builded := languagesBuildCommands[path.Ext(handler.File)]
				handlerExecutablePath := handler.File
				if builded {
					handlerExecutablePath = path.Join(handlersFilesPath, handlerName, "executable")
				}

				runCommand(
					handlerExecutablePath,
					writer,
					request,
					func(message string, err error) {
						message = fmt.Sprintf("[XServer] [%s Handler] [Error] %s: %s", handlerName, message, err)
						logger.Error(message)
						writer.Write([]byte(message + "\n"))
					},
					func(message string) {
						logger.Info(fmt.Sprintf("[XServer] [%s Handler] %s", handlerName, message))
					},
				)
			},
			"POST",
		)
	}

	err := xserver.Start()
	if err != nil {
		return err
	}
	return nil
}

func usage() {
	fmt.Println("usage: xserver <command>")
	fmt.Println("\tcommands:")
	fmt.Println("\t\tbuild: compiles all handlers")
	fmt.Println("\t\tstart: start server")
}

func main() {
	arguments := os.Args
	if (len(arguments)) == 1 {
		usage()
		return
	}
	command, ok := commands[arguments[1]]
	if !ok {
		usage()
		return
	}

	config, err := config.Load(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := logger.Configure(config); err != nil {
		fmt.Println(err)
		return
	}

	if err := command(config); err != nil {
		logger.Error(err.Error())
		return
	}

}
