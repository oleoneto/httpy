# httpy
HTTPy, a CLI tool for programmatically managing collections of HTTP requests.

**Table of Contents**
- [httpy](#httpy)
  - [Installation](#installation)
  - [Commands](#commands)
    - [httpy](#httpy)
    - [server](#server)
    - [http](#http)
    - [version](#version)
  - [To Do](#to-do)

## Installation
```
make install
```

## Commands
Assume the installed binary is called `httpy`.

### httpy
```
HTTPy, a CLI tool for programmatically managing collections of HTTP requests

Usage:
  httpy [flags]
  httpy [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  http        Make HTTP requests
  server      Run a mock HTTP server
  version     Shows the version of the CLI

Flags:
      --config-dir string   config directory (default "$HOME/.httpy")
      --db-url string       database url (default "httpy.sqlite3")
  -h, --help                help for httpy
  -o, --output string       output format (default "yaml")
      --time                time executions
      --verbose             enable detailed logging
```

### server
```
Run a mock HTTP server

Usage:
  httpy server [flags]

Flags:
  -f, --file string
  -h, --help          help for server
  -p, --port int       (default 3333)
  -r, --show-routes
```

### http
```
Make HTTP requests

Usage:
  httpy http [flags]

Flags:
  -f, --file string
  -h, --help          help for http
```

### version
```
Shows the version of the CLI

Usage:
  httpy version [flags]

Flags:
  -h, --help   help for version
```

## To Do
[Check out open issues](https://github.com/oleoneto/httpy/issues).
