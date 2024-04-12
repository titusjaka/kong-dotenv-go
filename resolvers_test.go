package kongdotenv_test

import (
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"

	kongdotenv "github.com/titusjaka/kong-dotenv-go"
)

func TestENVFileReader(t *testing.T) {
	t.Parallel()

	t.Run("Basic parser", func(t *testing.T) {
		t.Parallel()

		var cli struct {
			String string `env:"STRING"`
			Int    int    `env:"INT"`
			Bool   bool   `env:"BOOL"`
		}

		envFile := `STRING=🍕
INT=5
BOOL=true
`

		r, err := kongdotenv.ENVFileReader(strings.NewReader(envFile))
		require.NoError(t, err)

		parser := mustNew(t, &cli, kong.Resolvers(r))
		_, err = parser.Parse([]string{})
		require.NoError(t, err)
		require.Equal(t, "🍕", cli.String)
		require.Equal(t, 5, cli.Int)
		require.True(t, cli.Bool)
	})

	t.Run("Substitutions", func(t *testing.T) {
		t.Parallel()

		var cli struct {
			String  string `env:"STRING"`
			Int     int    `env:"INT"`
			Bool    bool   `env:"BOOL"`
			String2 string `env:"STRING_2"`
		}

		envFile := `STRING=🍕
INT=5
BOOL=true
STRING_2=$STRING
`

		r, err := kongdotenv.ENVFileReader(strings.NewReader(envFile))
		require.NoError(t, err)

		parser := mustNew(t, &cli, kong.Resolvers(r))
		_, err = parser.Parse([]string{})
		require.NoError(t, err)
		require.Equal(t, "🍕", cli.String)
		require.Equal(t, 5, cli.Int)
		require.True(t, cli.Bool)
	})

	t.Run("Prioritize env over envfile", func(t *testing.T) {
		t.Parallel()

		defer func() {
			_ = os.Unsetenv("PIZZA_STRING")
		}()

		require.NoError(t, os.Setenv("PIZZA_STRING", "pizza"))

		var cli struct {
			String string `kong:"env=PIZZA_STRING"`
		}

		envFile := `PIZZA_STRING=🍕`

		r, err := kongdotenv.ENVFileReader(strings.NewReader(envFile))
		require.NoError(t, err)

		parser := mustNew(t, &cli, kong.Resolvers(r))
		_, err = parser.Parse([]string{})
		require.NoError(t, err)
		require.Equal(t, "pizza", cli.String)
	})

	t.Run("Prioritize envfile over default", func(t *testing.T) {
		t.Parallel()

		defer func() {
			_ = os.Unsetenv("PIZZA_STRING1")
		}()

		require.NoError(t, os.Setenv("PIZZA_STRING1", "pizza"))

		var cli struct {
			String string `kong:"env=PIZZA_STRING1,default='pepperoni'"`
		}

		envFile := `PIZZA_STRING1=🍕`

		r, err := kongdotenv.ENVFileReader(strings.NewReader(envFile))
		require.NoError(t, err)

		parser := mustNew(t, &cli, kong.Resolvers(r))
		_, err = parser.Parse([]string{})
		require.NoError(t, err)
		require.Equal(t, "pizza", cli.String)
	})

	t.Run("Check multiple ENVs", func(t *testing.T) {
		t.Parallel()

		var cli struct {
			String string `env:"STRING,STRING2"`
		}

		envFile := `STRING2=🍕
STRING=pizza
`

		r, err := kongdotenv.ENVFileReader(strings.NewReader(envFile))
		require.NoError(t, err)

		parser := mustNew(t, &cli, kong.Resolvers(r))
		_, err = parser.Parse([]string{})
		require.NoError(t, err)
		require.Equal(t, "pizza", cli.String)
	})
}

func mustNew(t *testing.T, cli interface{}, options ...kong.Option) *kong.Kong {
	t.Helper()
	options = append([]kong.Option{
		kong.Name("test"),
		kong.Exit(func(int) {
			t.Helper()
			t.Fatalf("unexpected exit()")
		}),
	}, options...)
	parser, err := kong.New(cli, options...)
	require.NoError(t, err)
	return parser
}
