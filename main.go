package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"xserver/src/builders"
	"xserver/src/config"
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
	handlersPath           = "bin/handlers/"
	languagesBuildCommands = map[string]func(string, string) error{
		".go": builders.Go,
	}
	languagesRunCommands = map[string]func(string, http.ResponseWriter, *http.Request, func(err error)){
		".go": runners.Executable,
	}
)

func create(config *config.Config) error {
	fmt.Println("XServer: create project.")
	return nil
}

func build(config *config.Config) error {
	fmt.Println("[XServer] [Build] Build project")

	if err := os.MkdirAll(handlersPath, os.ModePerm); err != nil {
		return fmt.Errorf("[XServer] [Build] [Error] failed create files directory: %s", err)
	}

	for handlerName, handler := range config.Handlers {
		fmt.Println(fmt.Sprintf(`[XServer] [Build] build "%s" handler`, handlerName))
		//file := path.Base(handler.Path)
		fileExtention := path.Ext(handler.Path)
		buildCommand, ok := languagesBuildCommands[fileExtention]
		if ok {
			if err := buildCommand(handler.Path, path.Join(handlersPath, handlerName, "executable")); err != nil {
				return fmt.Errorf(`[XServer] [Build] [Error] failed compile "%s" handler: %s`, handlerName, err)
			}
		}
	}
	return nil
}

func start(config *config.Config) error {
	fmt.Println("[XServer] Start project")
	xserver := server.Create(
		&server.ServerOptions{Url: config.Url},
	)

	for handlerName, handler := range config.Handlers {
		xserver.AddHandler(
			"/"+handlerName,
			func(writer http.ResponseWriter, request *http.Request) {
				runCommand, ok := languagesRunCommands[path.Ext(handler.Path)]
				if !ok {
					writer.Write([]byte(fmt.Sprintf("[XServer] [%s Handler] [Error] run command is unknown", handlerName)))
				}

				_, builded := languagesBuildCommands[path.Ext(handler.Path)]
				handlerPath := handler.Path
				if builded {
					handlerPath = path.Join(handlersPath, handlerName, "executable")
				}

				runCommand(
					handlerPath,
					writer,
					request,
					func(err error) {
						writer.Write([]byte(fmt.Sprintf("[XServer] [%s Handler] [Error] failed run command: %s", handlerName, err)))
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
	fmt.Println("usage: xserver <command> [command args ...]")
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
	if err := command(config); err != nil {
		fmt.Println(err)
		return
	}

}
