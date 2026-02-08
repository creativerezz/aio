package cli

import (
	"github.com/creativerezz/aio/internal/plugins/db/fsdb"
)

// handleManagementCommands handles management-related commands (delete, print, etc.)
// Returns (handled, error) where handled indicates if a command was processed and should exit
func handleManagementCommands(currentFlags *Flags, aioDb *fsdb.Db) (handled bool, err error) {
	if currentFlags.WipeContext != "" {
		err = aioDb.Contexts.Delete(currentFlags.WipeContext)
		return true, err
	}

	if currentFlags.WipeSession != "" {
		err = aioDb.Sessions.Delete(currentFlags.WipeSession)
		return true, err
	}

	if currentFlags.PrintSession != "" {
		err = aioDb.Sessions.PrintSession(currentFlags.PrintSession)
		return true, err
	}

	if currentFlags.PrintContext != "" {
		err = aioDb.Contexts.PrintContext(currentFlags.PrintContext)
		return true, err
	}

	return false, nil
}
