# XServer
Server with multiple programming languages support.
___
## Config
Create `config.yml` file to define handlers.
```yaml
---
url: ip:port
logs: path/to/logs.txt

handlers:
  handler_name:
    path: /server/handler/path
    file: path/to/handler.ext
    build:
        tool: tool_for_build
        flags:
            - build_flag_1
            - build_flag_2
            ...
            - build_flag_n
    run:
        tool: tool_for_run
        arguments:
            - argument_1
            - argument_2
            ...
            - argument_n
```

## Build project
If your project contains compiled programming languages, you need to build all handlers of this type before start.

`$ xserver build`

All collected files will be placed in the `bin` directory.
The path to each handler corresponds to `bin/handler_name/executable`.

## Start server
`$ xserver start`
