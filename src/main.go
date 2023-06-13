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
	"xserver/src/utils"
)

var (
	commands = map[string]func(config *config.Config) error{
		"build": build,
		"start": start,
	}
	configPath             = "./config.yml"
	handlersFilesPath      = "bin/handlers/"
	languagesBuildCommands = map[string]func(string, string, ...string) error{
		".go":  builders.Go,
		".cpp": builders.Cpp,
	}
	languagesRunCommands = map[string]func(string, http.ResponseWriter, *http.Request, func(string, error), func(string), ...string){
		".go":  runners.Executable,
		".cpp": runners.Executable,
	}
)

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
			if err := builders.Tool(handler.Build.Tool, handler.File, path.Join(handlersFilesPath, handlerName, "executable"), handler.Build.Flags...); err != nil {
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
		} else {
			utils.CopyFile(handler.File, path.Join(handlersFilesPath, handlerName, path.Base(handler.File)))
		}
	}
	return nil
}

func start(config *config.Config) error {
	logger.Info("[XServer] Start project")

	for handlerName, handler := range config.Handlers {
		currentHandlerName := handlerName
		currentHandler := handler

		_, builded := languagesBuildCommands[path.Ext(currentHandler.File)]
		builded = builded || currentHandler.Build != nil
		handlerExecutablePath := path.Join(handlersFilesPath, currentHandlerName, path.Base(currentHandler.File))
		if builded {
			handlerExecutablePath = path.Join(handlersFilesPath, currentHandlerName, "executable")
		}

		runCommand := languagesRunCommands[path.Ext(currentHandler.File)]

		if currentHandler.Run != nil && currentHandler.Run.Tool != "" {
			runCommand = func(path string, writer http.ResponseWriter, request *http.Request, errorCallback func(string, error), logCallback func(string), args ...string) {
				runners.Tool(currentHandler.Run.Tool, path, writer, request, errorCallback, logCallback, args...)
			}
		}

		if runCommand == nil {
			if builded {
				runCommand = runners.Executable
			} else {
				message := fmt.Sprintf("[XServer] [%s Handler] [Error] run command is unknown", currentHandlerName)
				logger.Error(message)
				continue
			}
		}

		args := []string{}
		if currentHandler.Run != nil {
			args = append(args, currentHandler.Run.Args...)
		}

		server.AddHandler(
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
					args...,
				)
			},
		)
	}

	err := server.Start(config)
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
