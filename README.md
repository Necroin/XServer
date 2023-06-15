# XServer
Server with multiple programming languages support.

[Project Example](/example/)
___
## Supported languages
At the moment, the standard languages for assembly are:
- `Golang`
- `C/C++`

For run:
- `Golang`
- `C/C++`
- `Python`
- `Lua`

If your language is not in the list of standard languages or you want to use additional options for the build and/or run, you can use the `build` and/or `run` options in the configuration file for handler or task.
___
## Assembly
Build manually for your platform
```shell
go build -o xserver src/main.go
```
___
## Configuration file
The configuration file uses the `yaml` format.

Server uses the following configuration file structure:
- `url` - server url
- `log` - path to log file
- `handlers` - section for server handlers
  - `handler name` - defines the handler and makes it unique
    - `path` - server handler path
    - `file` - path to handler file
    - `build` - use for custom build, optional
      - `tool` - tool for build e.g. `gcc`/`g++`, optional
      - `flags` -  list of build flags, optional
    - `run` - use for custom handler run, optional
      - `tool` - tool for run e.g. `python`/`lua`, optional
      - `flags` -  list of run flags, optional
- `tasks` - section for server tasks
  - `handler name` - defines the task and makes it unique
    - `file` - path to handler file
    - `period` - cron formatted period
    - `build` - same as in `handlers` section
    - `run` - same as in `handlers` section
___
## Usage
### 1. Create Config
Create `config.yml` file to define server, handlers and periodic tasks.

### 2. Build project
You need to build all handlers and tasks before starting server.
```shell
$ xserver build
```
All files specified in the part `handlers` or `tasks` will be placed in the `bin` directory.

### 3. Start server
```shell
$ xserver start
```
