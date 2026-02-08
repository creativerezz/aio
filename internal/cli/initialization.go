package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/creativerezz/aio/internal/core"
	"github.com/creativerezz/aio/internal/i18n"
	"github.com/creativerezz/aio/internal/plugins/db/fsdb"
)

const ConfigDirPerms os.FileMode = 0755
const EnvFilePerms os.FileMode = 0644

// initializeAio initializes the aio database and plugin registry
func initializeAio() (registry *core.PluginRegistry, err error) {
	var homedir string
	if homedir, err = os.UserHomeDir(); err != nil {
		return
	}

	aioDb := fsdb.NewDb(filepath.Join(homedir, ".config/aio"))
	if err = aioDb.Configure(); err != nil {
		return
	}

	if registry, err = core.NewPluginRegistry(aioDb); err != nil {
		return
	}

	return
}

// ensureEnvFile checks for the default ~/.config/aio/.env file and creates it
// along with the parent directory if it does not exist.
func ensureEnvFile() (err error) {
	var homedir string
	if homedir, err = os.UserHomeDir(); err != nil {
		return fmt.Errorf("%s", fmt.Sprintf(i18n.T("could_not_determine_home_dir"), err))
	}
	configDir := filepath.Join(homedir, ".config", "aio")
	envPath := filepath.Join(configDir, ".env")

	if _, statErr := os.Stat(envPath); statErr != nil {
		if !os.IsNotExist(statErr) {
			return fmt.Errorf("%s", fmt.Sprintf(i18n.T("could_not_stat_env_file"), statErr))
		}
		if err = os.MkdirAll(configDir, ConfigDirPerms); err != nil {
			return fmt.Errorf("%s", fmt.Sprintf(i18n.T("could_not_create_config_dir"), err))
		}
		if err = os.WriteFile(envPath, []byte{}, EnvFilePerms); err != nil {
			return fmt.Errorf("%s", fmt.Sprintf(i18n.T("could_not_create_env_file"), err))
		}
	}
	return
}
