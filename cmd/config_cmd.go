package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/david-truong/liferay-issues-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long:  "Set a configuration value. Keys use dot notation (e.g. jira.instance, auth.email, auth.token, jira.default_project, defaults.issue_type).",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		if err := config.EnsureConfigDir(); err != nil {
			return fmt.Errorf("creating config directory: %w", err)
		}

		viper.Set(key, value)

		if err := viper.WriteConfigAs(config.ConfigFilePath()); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}

		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := viper.Get(key)
		if value == nil {
			fmt.Fprintf(os.Stderr, "Key %q not set\n", key)
			return nil
		}
		fmt.Println(value)
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all config values",
	RunE: func(cmd *cobra.Command, args []string) error {
		settings := viper.AllSettings()
		lines := flattenMap("", settings)
		sort.Strings(lines)
		for _, line := range lines {
			fmt.Println(line)
		}
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show config file path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.ConfigFilePath())
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configPathCmd)
}

func flattenMap(prefix string, m map[string]interface{}) []string {
	var lines []string
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			lines = append(lines, flattenMap(key, val)...)
		default:
			s := fmt.Sprintf("%v", val)
			if s == "" {
				continue
			}
			// Mask token for security
			if strings.Contains(key, "token") {
				if len(s) > 4 {
					s = s[:4] + strings.Repeat("*", len(s)-4)
				}
			}
			lines = append(lines, fmt.Sprintf("%s = %s", key, s))
		}
	}
	return lines
}
