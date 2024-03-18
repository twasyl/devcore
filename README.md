`devcore` aims to help developers manage their environment. Discover how it can help you.

# Tools

## Install

`devcore` lets you install some tools that may be useful for you. These tools are downloaded and installed in your local
environment.

To install a tool:
```bash
$ devcore tools install <toolname>
# or
$ devcore tools install <toolname> -v <version>
```

## List

Listing available tools:
```bash
$ devcore tools list
# or
$ devcore tools list -v
```

# Servers

## Install

`devcore` lets you install some servers that may be useful for you. These servers are downloaded and installed in your local
environment.

To install a tool:
```bash
$ devcore servers install <servername>
# or
$ devcore servers install <servername> -v <version>
```

The server will be installed under `$(devcore config get servers-dir)`

## List

Listing available servers:
```bash
$ devcore servers list
# or
$ devcore servers list -v
```

# Managing docker compose contexts

`devcore` lets you manage docker compose contexts. Commands:
- `devcore docker-compose ...`
- `devcore dc ...`

## Create

```bash
$ devcore docker-compose create -n <name> -f <docker compose file> [-d <description>]
```

**Note:** this automatically set the newly created context as the current one.

## List

```bash
$ devcore docker-compose list [-v]
```

## Delete

```bash
$ devcore docker-compose delete [name]
```

If no name is provided, the current docker compose context is deleted.

## Set current (or default)

```bash
$ devcore docker-compose set-current <name>
```

## Start

```bash
$ devcore docker-compose start [-v] [<name>]
```

If no name is provided, the current docker compose context is started.

## Stop

```bash
$ devcore docker-compose stop [-v] [<name>]
```

If no name is provided, the current docker compose context is stopped.

# Open folder

Open the folder containing the docker compose file of the context.

```bash
$ devcore docker-compose open-folder [<name>]
```

If no name is provided, the current docker compose context is used.