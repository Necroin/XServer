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
	languagesBuildCommands = map[string]func(string, string, ...string) error{
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

	if err := os.RemoveAll(handlersFilesPath); err != nil {
		return fmt.Errorf("[XServer] [Build] [Error] failed delete handlers file directory: %s", err)
	}

	if err := os.MkdirAll(handlersFilesPath, os.ModePerm); err != nil {
		return fmt.Errorf("[XServer] [Build] [Error] failed create handlers file directory: %s", err)
	}

	for handlerName, handler := range config.Handlers {
		logger.Info(fmt.Sprintf(`[XServer] [Build] build "%s" handler`, handlerName))
		if err := os.MkdirAll(path.Join(handlersFilesPath, handlerName), os.ModePerm); err != nil {
			return fmt.Errorf("[XServer] [Build] [Error] failed create handler file directory: %s", err)
		}

		if handler.Build != nil && handler.Build.Tool != "" {
			logger.Info(fmt.Sprintf(`[XServer] [Build] handler "%s" has specified build options -> build by options`, handlerName))
			if err := builders.Custom(handler.Build.Tool, handler.File, path.Join(handlersFilesPath, handlerName, "executable"), handler.Build.Flags...); err != nil {
				logger.Error(fmt.Sprintf(`[XServer] [Build] [Error] failed compile "%s" handler: %s`, handlerName, err))
			}
			continue
		}

		buildCommand, ok := languagesBuildCommands[path.Ext(handler.File)]

		if ok {
			flags := []string{}
			if handler.Build != nil {
				flags = handler.Build.Flags
			}
			if err := buildCommand(handler.File, path.Join(handlersFilesPath, handlerName, "executable"), flags...); err != nil {
				logger.Error(fmt.Sprintf(`[XServer] [Build] [Error] failed compile "%s" handler: %s`, handlerName, err))
			}
			continue
		}
		logger.Info(fmt.Sprintf(`[XServer] [Build] handler "%s" has not any build options -> skip`, handlerName))
	}
	return nil
}

func start(config *config.Config) error {
	logger.Info("[XServer] Start project")
	xserver := server.Create(
		&server.ServerOptions{Url: config.Url},
	)

	for handlerName, handler := range config.Handlers {
		currentHandlerName := handlerName
		currentHandler := handler

		_, builded := languagesBuildCommands[path.Ext(currentHandler.File)]
		builded = builded || handler.Build != nil
		handlerExecutablePath := currentHandler.File
		if builded {
			handlerExecutablePath = path.Join(handlersFilesPath, currentHandlerName, "executable")
		}

		runCommand, ok := languagesRunCommands[path.Ext(currentHandler.File)]

		if !ok {
			if builded {
				runCommand = runners.Executable
			} else {
				message := fmt.Sprintf("[XServer] [%s Handler] [Error] run command is unknown", currentHandlerName)
				logger.Error(message)
			}
		}

		if handler.Run != nil && handler.Run.Tool != "" {
			runCommand = func(path string, writer http.ResponseWriter, request *http.Request, errorCallback func(string, error), logCallback func(string)) {
				runners.Tool(currentHandler.Run.Tool, path, writer, request, errorCallback, logCallback)
			}
		}

		xserver.AddHandler(
			currentHandler.Path,
			func(writer http.ResponseWriter, request *http.Request) {
				logger.Info(fmt.Sprintf("[XServer] [%s Handler] handler called", currentHandlerName))

				runCommand(
					handlerExecutablePath,
					writer,
					request,
					func(message string, err error) {
						message = fmt.Sprintf("[XServer] [%s Handler] [Error] %s: %s", currentHandlerName, message, err)
						logger.Error(message)
						writer.Write([]byte(message + "\n"))
					},
					func(message string) {
						logger.Info(fmt.Sprintf("[XServer] [%s Handler] %s", currentHandlerName, message))
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
