package cli

import (
	"github.com/spf13/cobra"
	"github.com/testy/lightf/pkg/slog"
)

func runContextCmd(cmd *cobra.Command, args []string) {
	cfg := contextConfig // local copy of pointer, from cli.go file
	provider := configProvider

	if !cfg.HasAnyContext() {
		slog.Fatalf("There is no context available. Maybe wrong configuration?")
	}

	if len(args) == 0 {
		// display all contexts
		slog.Infoln("Available contexts:")
		for _, serv := range cfg.Servers {
			slog.Infof("    - %v [%q]", serv.Name, serv.Address)
		}

		_, err := cfg.CurrentContext()
		var suffix string
		if err != nil {
			suffix = " (invalid)"
		}
		slog.Infof("Current context: %v%v", cfg.CurrentServer, suffix)
		return
	}

	// switch to given context
	arg := args[0]
	if !cfg.HasContext(arg) {
		slog.Fatalf("Given context does not exist, ctx=%q", arg)
	}
	res := cfg.SetCurrentServer(arg)
	if !res {
		slog.Infof("No changes, already in this context, ctx=%q", arg)
	} else {
		slog.Infof("Switched to context %q.", arg)

		// HARDCODE, we could pull this out of here, but really not neccessary.
		provider.Set("current-server", arg)
		provider.WriteConfig()
	}
}
