# XServer
Server with multiple programming languages support.
___
## Create Config
Create `config.yml` file to define server parameters and handlers.

https://github.com/Necroin/XServer/blob/ef4efe846259a1a0c2f65efd97a462a20df6efbe/example/config.yml#L1-L33

## Build project
If your project contains compiled programming languages, you need to build all handlers of this type before start.

`$ xserver build`

All files specified in the part `handlers` will be placed in the `bin` directory.

The path to each handler corresponds to `bin/handlers/handler_name/executable` if the handler requires assembly.
Or `bin/handlers/handler_name/handler_file_name.handler_file_extension` if the handler is interpreted.

## Start server
`$ xserver start`
