# XServer
Server with multiple programming languages support.

[Project Example](/example/)
___
## Supported languages
At the moment, the standard languages for assembly and run are:
- `Golang`
- `C/C++`

If your language is not in the list of standard languages or you want to use additional options for the build and/or, you can use the `build` and/or `run` options in the configuration file for handler or task.
___
## Assembly
Build manually for your platform
```shell
go build -o xserver src/main.go
```
___
## Usage
### 1. Create Config
Create `config.yml` file to define server, handlers and periodic tasks.

https://github.com/Necroin/XServer/blob/ef4efe846259a1a0c2f65efd97a462a20df6efbe/example/config.yml#L1-L33

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
