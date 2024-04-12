# ENV File Resolver for [kong](https://github.com/alecthomas/kong)

This tiny package adds ENV-files loader to the [kong](https://github.com/alecthomas/kong) CLI parser.

## Usage

The easiest way to add ENV-file support to the app is the following.

```dotenv
# ./.env
VAR_STRING="Any value here"
```

```go
// ./main.go

package main

import (
	"fmt"

	"github.com/alecthomas/kong"

	kongdotenv "github.com/titusjaka/kong-dotenv-go"
)

type CLI struct {
	EnvFile kongdotenv.ENVFileConfig `kong:"optional,default='.env',name=env-file"`

	EnvVarString string `kong:"optional,name=var-string,env=VAR_STRING"`
}

func (c CLI) Run() error {
	fmt.Printf(`Env File: %s; Env Var String: %s.`, c.EnvFile, c.EnvVarString,
	)

	return nil
}

func main() {
	var cli CLI

	kCtx := kong.Parse(&cli)
	kCtx.FatalIfErrorf(kCtx.Run())
}
```

```shell
go run ./main.go --env-file=./.env

# Output: Env File: .env; Env Var String: Any value here.
```

## Limits

The package has the following limitations:
- Values are applied only to variables with an env tag.
- ENV variables have higher priority than ENV-file ones. Therefore, if you have an ENV set, the same variable in the ENV file will be ignored.

## Examples

- example with minimum configuration: [./examples/config_type](./examples/config_type);
- obsolete example with `kong.ConfigurationLoader` resolver: [./examples/resolver](./examples/resolver).
