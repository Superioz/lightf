package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/testy/lightf/internal/env"
	"github.com/testy/lightf/internal/flags"
	"github.com/testy/lightf/pkg/slog"
	"github.com/testy/lightf/pkg/token"
)

func runTokenCmd(cmd *cobra.Command, args []string) {
	length, err := cmd.Flags().GetInt(flags.Length)
	if err != nil {
		slog.Fatal(err)
	}
	charset, err := cmd.Flags().GetString(flags.Charset)
	if err != nil {
		slog.Fatal(err)
	}

	gen := token.G(charset, length, &token.Settings{
		MinTokenLength:   viper.GetInt(env.TokenGenMinLength),
		MinCharsetLength: viper.GetInt(env.TokenGenMinCharsetLength),
		MinExpireLength:  viper.GetInt64(env.TokenGenMinExpireLength),
	})

	tok, err := gen.Do()
	if err != nil {
		slog.Fatal(err)
	}

	slog.Infof("Token: %q", tok.Str)
	slog.Infof("Timestamp: %v", tok.Creation/1000000)
}
