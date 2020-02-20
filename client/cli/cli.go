// package to create and execute the lightf cli tool.
// this has nothing to do with the server managing.
package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/testy/lightf/client/context"
	"github.com/testy/lightf/internal/env"
	"github.com/testy/lightf/internal/flags"
	"github.com/testy/lightf/internal/version"
	"github.com/testy/lightf/pkg/slog"
)

var (
	// the root cmd, where nothing really happens
	// this is just to get the command help
	// and to access the command's children.
	//
	// Usage: lightf [flags]
	rootCmd = &cobra.Command{
		Use:   "lightf",
		Short: "Lightweight, secure and transient fileserver written in Go.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// get global flag "verbose" for logging
			verbose, err := cmd.Flags().GetBool(flags.Verbose)
			if err == nil && verbose {
				slog.SetVerbosity(1)
			}
		},
	}

	// version command, to simply display
	// the current version of this executable (the client)
	//
	// Usage: lightf version [--short]
	versionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "Displays the current cli version",
		Run: func(cmd *cobra.Command, args []string) {
			short, err := cmd.Flags().GetBool(flags.Short)
			if err != nil {
				slog.Fatal(err)
			}

			if short {
				slog.Infof(version.ClientVersion)
			} else {
				slog.Infof("You are using lightf version %q", version.ClientVersion)
			}
		},
	}

	// token command, to generate a token
	// used for authentication with a server.
	// This command can be helpful, if you are operating
	// a server and need a valid token.
	//
	// Usage: lightf token [--length <length>] [--charset <charset>]
	tokenCmd = &cobra.Command{
		Use:   "token",
		Short: "Generates a new valid token",
		Long: `Generates a new valid lightf token with supplied options.
If no options are supplied, secure defaults will be used.
Tokens can be a potentially security risk when not using them correctly.`,
		Run: runTokenCmd,
	}

	// context command, to switch to a specific
	// context or just display one.
	contextCmd = &cobra.Command{
		Use:     "context [context]",
		Aliases: []string{"ctx"},
		Short:   "Switches or displays available contexts",
		Long: `When no arguments given, simply displays all available contexts.
These contexts represent a server, configured in a context.yml.
Otherwise use the argument to switch to the given context.`,
		Run: runContextCmd,
	}

	uploadCmd = &cobra.Command{
		Use:     "upload <file>",
		Aliases: []string{"up"},
		Short:   "Uploads a file to a lightf server",
		Run:     runUploadCmd,
	}
)

var (
	contextConfig *context.Config

	// we need to hold a second instance of viper
	// as we otherwise get in trouble with using
	// env variables and a config at the same time.
	configProvider *viper.Viper
)

func Exec() error {
	return rootCmd.Execute()
}

func init() {
	// ================================
	// get viper to set default envs
	// ================================
	viper.AutomaticEnv()

	viper.SetDefault(env.VerboseLogging, false)

	viper.SetDefault(env.VersionShortened, false)

	viper.SetDefault(env.TokenGenMinLength, 3)
	viper.SetDefault(env.TokenGenMinCharsetLength, 3)
	viper.SetDefault(env.TokenGenMinExpireLength, int64(10))
	viper.SetDefault(env.TokenGenDefaultLength, 18)
	viper.SetDefault(env.TokenGenDefaultCharset, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	viper.SetDefault(env.ContextConfigPath, "$HOME/.lightf")

	viper.SetDefault(env.UploadTimeoutInSeconds, 10)

	// ================================
	// get the context.yml config
	// ================================
	configProvider = viper.New()
	configProvider.SetConfigName("context")
	configProvider.SetConfigType("yml")

	configProvider.AddConfigPath(viper.GetString(env.ContextConfigPath))
	configProvider.AddConfigPath(".")

	if err := configProvider.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
		// actually we can ignore both cases.
		// if no config is set, the user will run into a message
		// telling him, that no context could be found, and that he should
		// use the context flags. No error.
	}
	configProvider.Unmarshal(&contextConfig)

	// ================================
	// setup the commands
	// ================================
	rootCmd.AddCommand(versionCmd, tokenCmd, contextCmd, uploadCmd)

	rootCmd.PersistentFlags().BoolP(flags.Verbose, flags.V, viper.GetBool(env.VerboseLogging), "verbose logging")

	versionCmd.Flags().BoolP(flags.Short, flags.S, viper.GetBool(env.VersionShortened), "only version")

	tokenCmd.Flags().IntP(flags.Length, flags.L, viper.GetInt(env.TokenGenDefaultLength), "token length")
	tokenCmd.Flags().StringP(flags.Charset, flags.C, viper.GetString(env.TokenGenDefaultCharset), "charset")

	uploadCmd.Flags().StringP(flags.Server, flags.S, "", "server address")
	uploadCmd.Flags().StringP(flags.Context, flags.C, "", "context")
	uploadCmd.Flags().StringP(flags.Token, flags.T, "", "auth token")
	uploadCmd.Flags().Int(flags.Timeout, viper.GetInt(env.UploadTimeoutInSeconds), "timeout in seconds")
}
