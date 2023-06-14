package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"xserver/src/builders"
	"xserver/src/config"
	"xserver/src/logger"
	"xserver/src/runners"
	"xserver/src/server"
	"xserver/src/utils"

	"github.com/robfig/cron"
)

var (
	commands = map[string]func(config *config.Config) error{
		"build": build,
		"start": start,
	}
	configPath        = "./config.yml"
	handlersFilesPath = "bin/handlers/"
	tasksFilesPath    = "bin/tasks/"

	languagesBuildCommands = map[string]func(string, string, ...string) error{
		".go":  builders.Go,
		".cpp": builders.Cpp,
	}
	languagesRunCommands = map[string]func(string, io.Writer, io.Reader, func(string, error), func(string), ...string){
		".go":  runners.Executable,
		".cpp": runners.Executable,
	}
)

func buildUnits(unitTag string, unitsFilesPath string, units map[string]config.ExecutableServerUnit) error {
	if err := os.RemoveAll(unitsFilesPath); err != nil {
		return fmt.Errorf("[XServer] [Build] [%s] [Error] failed delete file directory: %s", unitTag, err)
	}

	if err := os.MkdirAll(unitsFilesPath, os.ModePerm); err != nil {
		return fmt.Errorf("[XServer] [Build] [%s] [Error] failed create file directory: %s", unitTag, err)
	}

	for unitName, unit := range units {
		logger.Info(fmt.Sprintf(`[XServer] [Build] [%s] build "%s"`, unitTag, unitName))
		if err := os.MkdirAll(path.Join(unitsFilesPath, unitName), os.ModePerm); err != nil {
			return fmt.Errorf("[XServer] [Build] [%s] [Error] failed create file directory: %s", unitTag, err)
		}

		if unit.Build != nil && unit.Build.Tool != "" {
			logger.Info(fmt.Sprintf(`[XServer] [Build] [%s] "%s" has specified build options -> build by options`, unitTag, unitName))
			if err := builders.Tool(unit.Build.Tool, unit.File, path.Join(unitsFilesPath, unitName, "executable"), unit.Build.Flags...); err != nil {
				logger.Error(fmt.Sprintf(`[XServer] [Build] [%s] [Error] failed compile "%s": %s`, unitTag, unitName, err))
			}
			continue
		}

		buildCommand, ok := languagesBuildCommands[path.Ext(unit.File)]

		if ok {
			flags := []string{}
			if unit.Build != nil {
				flags = unit.Build.Flags
			}
			if err := buildCommand(unit.File, path.Join(unitsFilesPath, unitName, "executable"), flags...); err != nil {
				logger.Error(fmt.Sprintf(`[XServer] [Build] [%s] [Error] failed compile "%s": %s`, unitTag, unitName, err))
			}
			continue
		} else {
			if err := utils.CopyFile(unit.File, path.Join(unitsFilesPath, unitName, path.Base(unit.File))); err != nil {
				logger.Error(fmt.Sprintf(`[XServer] [Build] [%s] [Error] failed copy "%s": %s`, unitTag, unitName, err))
			}
		}
	}

	return nil
}

func build(config *config.Config) error {
	logger.Info("[XServer] [Build] Build project")

	if err := buildUnits("Handlers", handlersFilesPath, config.Handlers); err != nil {
		return err
	}

	if err := buildUnits("Tasks", tasksFilesPath, config.Tasks); err != nil {
		return err
	}

	return nil
}

func getUnitRunCommand(unitTag string, unitsFilesPath string, unitName string, unit config.ExecutableServerUnit) (func(io.Writer, io.Reader), error) {
	_, stdBuilded := languagesBuildCommands[path.Ext(unit.File)]
	builded := stdBuilded || (unit.Build != nil)
	unitExecutablePath := path.Join(unitsFilesPath, unitName, path.Base(unit.File))
	if builded {
		unitExecutablePath = path.Join(unitsFilesPath, unitName, "executable")
	}

	runCommand := languagesRunCommands[path.Ext(unit.File)]

	if unit.Run != nil && unit.Run.Tool != "" {
		runCommand = func(path string, writer io.Writer, request io.Reader, errorCallback func(string, error), logCallback func(string), args ...string) {
			runners.Tool(unit.Run.Tool, path, writer, request, errorCallback, logCallback, args...)
		}
	}

	if runCommand == nil {
		if builded {
			runCommand = runners.Executable
		} else {
			return nil, fmt.Errorf(fmt.Sprintf("[XServer] [%s %s] [Error] run command is unknown", unitName, unitTag))
		}
	}

	args := []string{}
	if unit.Run != nil {
		args = unit.Run.Args
	}

	return func(writer io.Writer, request io.Reader) {
		runCommand(
			unitExecutablePath,
			writer,
			request,
			func(message string, err error) {
				message = fmt.Sprintf(`{ "error": "[XServer] [%s %s] [Error] %s: %s" }`, unitName, unitTag, message, err)
				logger.Error(message)
				writer.Write([]byte(message + "\n"))
			},
			func(message string) {
				logger.Info(fmt.Sprintf("[XServer] [%s %s] %s", unitName, unitTag, message))
			},
			args...,
		)
	}, nil
}

func start(config *config.Config) error {
	logger.Info("[XServer] Start project")

	for handlerName, handler := range config.Handlers {
		currentHandlerName := handlerName
		currentHandler := handler

		runCommand, err := getUnitRunCommand("Handler", handlersFilesPath, currentHandlerName, currentHandler)

		if err != nil {
			logger.Error(err.Error())
			continue
		}

		server.AddHandler(
			currentHandler.Path,
			func(writer http.ResponseWriter, request *http.Request) {
				logger.Info(fmt.Sprintf("[XServer] [%s Handler] handler called", currentHandlerName))
				runCommand(writer, request.Body)
			},
		)
	}

	cron := cron.New()
	for taskName, task := range config.Tasks {
		currentTaskName := taskName
		currentTask := task

		runCommand, err := getUnitRunCommand("Task", tasksFilesPath, currentTaskName, currentTask)

		if err != nil {
			logger.Error(err.Error())
			continue
		}

		cron.AddFunc(
			currentTask.Period,
			func() {
				logger.Info(fmt.Sprintf("[XServer] [%s Task] task started", currentTaskName))
				outBuffer := &bytes.Buffer{}
				runCommand(outBuffer, &bytes.Buffer{})
				logger.Info(fmt.Sprintf("[XServer] [%s Task] returned: %s", currentTaskName, outBuffer.String()))
			},
		)
	}
	cron.Start()

	err := server.Start(config)
	if err != nil {
		cron.Stop()
		return err
	}
	return nil
}

func usage() {
	fmt.Println("usage: xserver <command>")
	fmt.Println("\tcommands:")
	fmt.Println("\t\tbuild: compiles all handlers and tasks")
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
