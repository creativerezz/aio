package cli

import (
	"github.com/creativerezz/aio/internal/core"
	restapi "github.com/creativerezz/aio/internal/server"
)

// handleSetupAndServerCommands handles setup and server-related commands
// Returns (handled, error) where handled indicates if a command was processed and should exit
func handleSetupAndServerCommands(currentFlags *Flags, registry *core.PluginRegistry, version string) (handled bool, err error) {
	// if the setup flag is set, run the setup function
	if currentFlags.Setup {
		err = registry.Setup()
		return true, err
	}

	if currentFlags.Serve {
		registry.ConfigureVendors()
		err = restapi.Serve(registry, currentFlags.ServeAddress, currentFlags.ServeAPIKey)
		return true, err
	}

	if currentFlags.ServeOllama {
		registry.ConfigureVendors()
		err = restapi.ServeOllama(registry, currentFlags.ServeAddress, version)
		return true, err
	}

	return false, nil
}
