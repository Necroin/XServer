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
- `log` - path to log file (use `stdout` by default)
- `log_level` - `error`/`info`/`debug`/`verbose` (`info` by default)
- `database` - database options (`sqlite`)
  - `enable` - use database flag (`true`/`false`)
  - `storage` - path to storege `.db` file (`storage.db` by default)
  - `schema` - path to schema `.json` file (`schema.json` by default)
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
___
## Database
Server use sqlite database.

The scheme and operations uses the json format.
___
### Schema format
```
[
  {
    "name": "table_name",
    "fields": [
      {
        "name": "field_name",
        "type": "field_type",
        "nullable": true/false
      },
      ...
    ],
  "primary_key": ["field_name1", "field_name2", ...]
  },
  ...
]
```
___
### Operations
Database operations are implemented via server endpoints.
- `insert` - `/db/insert`
- `select` - `/db/select`
- `update` - `/db/update`
- `delete` - `/db/delete`
___
### Operations request format
- `insert`
```
{
  "table": "Users",
  "fields": [
    {
      "name": "field_name",
      "value": "field_value"
    },
    ...
  ]
}
```

- `select`
```
{
  "table": "Users",
  "fields": [{"name": "field_name"}, ...]
}
```

- `update`
```
{
  "table": "Users",
  "filters": [
      {
          "name": "field_name",
          "operator": "any_sql_compare_operator",
          "value": "field_value"
      },
      ...
  ],
  "fields": [
      {
          "name": "update_field_name",
          "value": "update_field_value"
      },
      ...
  ]
}
```

- `delete`
```
{
  "table": "Users",
  "filters": [
      {
          "name": "field_name",
          "operator": "any_sql_compare_operator",
          "value": "field_value"
      },
      ...
  ]
}
```