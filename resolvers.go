package kongdotenv

import (
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

// ENVFileConfig adds ENVFileReader configuration loader to load configuration
// from a .env file specified by a flag.
//
// Use this as a flag value to support loading of custom configuration via a flag.
type ENVFileConfig string

// BeforeResolve adds a resolver for .env file
func (c ENVFileConfig) BeforeResolve(ctx *kong.Context, trace *kong.Path) error {
	flag := ctx.FlagValue(trace.Flag)
	envFlag, ok := flag.(ENVFileConfig)
	if !ok {
		return fmt.Errorf("invalid type: '%T' expected, got '%T'", envFlag, flag)
	}

	return kong.Configuration(ENVFileReader, string(envFlag)).Apply(ctx.Kong)
}

// ENVFileReader returns a kong.Resolver that retrieves values from a .env file source.
//
// ENVFileReader resolves only flags with `env:"X"` tag.
func ENVFileReader(r io.Reader) (kong.Resolver, error) {
	values, err := godotenv.Parse(r)
	if err != nil {
		return nil, err
	}

	var f kong.ResolverFunc = func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
		// Skip, if flag doesn't have an environment variable.
		if len(flag.Envs) == 0 {
			return nil, nil
		}

		// Skip, if environment variable is already set.
		for _, env := range flag.Envs {
			if os.Getenv(env) != "" {
				return nil, nil
			}
		}

		for _, env := range flag.Envs {
			if raw, ok := values[env]; ok {
				return raw, nil
			}
		}

		return nil, nil
	}

	return f, nil
}
