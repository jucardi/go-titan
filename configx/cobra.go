package configx

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	RemoteEnvVariable = "CONFIG_REMOTE_URL"
	defaultSearchPath = "fixtures/config.yml"
	remoteFlag        = "remote"
	configFlag        = "config"
)

// LoadFlagsToCommand adds the configurator flags to a cobra command.
func LoadFlagsToCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(configFlag, "c", defaultSearchPath, "Use a local configuration file")
	cmd.PersistentFlags().StringP(remoteFlag, "r", "", "Use remote configuration")
}

// FromCommand loads the configuration based on the command flags received.
func FromCommand(cmd *cobra.Command) error {
	remote, _ := cmd.Flags().GetString(remoteFlag)
	filepath, _ := cmd.Flags().GetString(configFlag)
	if remote == "" {
		remote = os.Getenv(RemoteEnvVariable)
	}
	if remote != "" {
		log().Info("Loading config from remote: ", remote)
		return FromRemote(remote)
	}

	log().Info("Loading config from file")
	return FromFile(filepath)
}
