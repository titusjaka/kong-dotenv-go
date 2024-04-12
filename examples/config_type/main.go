package main

import (
	"fmt"

	"github.com/alecthomas/kong"

	kongdotenv "github.com/titusjaka/kong-dotenv-go"
)

type CLI struct {
	EnvFile kongdotenv.ENVFileConfig `kong:"optional,default='.env',name=env-file,help='Path to .env file'"`

	EnvVarString string `kong:"optional,name=env-var-string,env=ENV_VAR_STRING,help='String ENV variable'"`
	EnvVarInt    int    `kong:"optional,name=env-var-int,env=ENV_VAR_INT,help='Env variable #2'"`
}

func (c CLI) Run() error {
	fmt.Printf("Env File: %s; Env Var String: %s; Env Var Int: %d.\n",
		c.EnvFile,
		c.EnvVarString,
		c.EnvVarInt,
	)

	// Output: Env File: .env; Env Var String: ENVFileConfig; Env Var Int: 132.

	return nil
}

func main() {
	var cli CLI

	kCtx := kong.Parse(&cli)
	kCtx.FatalIfErrorf(kCtx.Run())
}
