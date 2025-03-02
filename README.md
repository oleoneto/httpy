# mock-http
Mock HTTP is a CLI tool for making requests and mocking server responses

**Table of Contents**
- [mock-http](#mock-http)
  - [Installation](#installation)
  - [Commands](#commands)
    - [mockhttp](#mockhttp)
    - [server](#server)
    - [http](#http)
    - [version](#version)
  - [To Do](#to-do)

## Installation
```
make install
```

## Commands
Assume the installed binary is called `mockhttp`.

### mockhttp
```
Usage:
  mockhttp [flags]
  mockhttp [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  http        Make HTTP requests
  server      Run a mock HTTP server
  version     Shows the version of the CLI

Flags:
  -h, --help                help for mockhttp
  -o, --output string       output format (default "yaml")
      --time                time executions
      --verbose             enable detailed logging

Use "mockhttp [command] --help" for more information about a command.
```

### server
```
Run a mock HTTP server

Usage:
  mockhttp server [flags]

Flags:
  -f, --file string
  -h, --help          help for server
  -p, --port int       (default 3333)
  -r, --show-routes

Global Flags:
  -o, --output string       output format (default "yaml")
      --time                time executions
      --verbose             enable detailed logging
```

### http
```
Make HTTP requests

Usage:
  mockhttp http [flags]

Flags:
  -f, --file string
  -h, --help          help for http

Global Flags:
  -o, --output string       output format (default "yaml")
      --time                time executions
      --verbose             enable detailed logging
```

### version
```
Shows the version of the CLI

Usage:
  mockhttp version [flags]

Flags:
  -h, --help   help for version

Global Flags:
  -o, --output string       output format (default "yaml")
      --time                time executions
      --verbose             enable detailed logging
```

## To Do
[Check out open issues](https://github.com/oleoneto/mock-http/issues).
