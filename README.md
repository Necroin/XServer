# XServer
Server with multiple programming languages support.
___
## Create Config
Create `config.yml` file to define server parameters and handlers.

https://github.com/Necroin/XServer/blob/b6a49e6d598e351473c545d2094763a12db82cb5/example/config.yml?plain=1#L1-L22

## Build project
If your project contains compiled programming languages, you need to build all handlers of this type before start.

`$ xserver build`

All files specified in the part `handlers` will be placed in the `bin` directory.

The path to each handler corresponds to `bin/handlers/handler_name/executable` if the handler requires assembly.
Or `bin/handlers/handler_name/handler_file_name.handler_file_extension` if the handler is interpreted.

## Start server
`$ xserver start`
