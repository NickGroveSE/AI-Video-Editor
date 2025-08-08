package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `Configure AI Video Editor settings including API keys, default parameters,
and processing preferences. Configuration is stored in ~/.ai-editor.yaml`,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set configuration values such as API keys and default processing parameters.

Available configuration keys:
  api-key          OpenAI API key for AI analysis
  whisper-model    Whisper model size (tiny, base, small, medium, large)
  default-duration Default clip duration
  default-quality  Default output quality (low, medium, high)
  temp-dir         Temporary directory for processing`,
	Args: cobra.ExactArgs(2),
	Example: `  # Set OpenAI API key
  ai-editor config set api-key sk-your-openai-key-here
  
  # Set default clip duration
  ai-editor config set default-duration 45s
  
  # Set Whisper model size
  ai-editor config set whisper-model medium`,
	RunE: runConfigSet,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  `Retrieve a configuration value by key.`,
	Args:  cobra.ExactArgs(1),
	Example: `  # Get current API key
  ai-editor config get api-key
  
  # Get default duration
  ai-editor config get default-duration`,
	RunE: runConfigGet,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Long:  `Display all current configuration settings.`,
	RunE:  runConfigList,
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Long:  `Reset all configuration values to their defaults. This will remove your custom settings.`,
	RunE:  runConfigReset,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configResetCmd)
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// Validate key
	if err := validateConfigKey(key); err != nil {
		return err
	}

	// Set the value
	viper.Set(key, value)

	// Ensure config directory exists
	configDir := filepath.Dir(viper.ConfigFileUsed())
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = home
		viper.SetConfigFile(filepath.Join(configDir, ".ai-editor.yaml"))
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config
	if err := viper.WriteConfig(); err != nil {
		// If config doesn't exist, create it
		if err := viper.SafeWriteConfig(); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	fmt.Printf("✅ Configuration updated: %s = %s\n", key, value)
	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := viper.GetString(key)
	
	if value == "" {
		fmt.Printf("❌ Configuration key '%s' not found or empty\n", key)
		return nil
	}

	// Hide sensitive values
	if key == "api-key" {
		if len(value) > 8 {
			value = value[:4] + "..." + value[len(value)-4:]
		}
	}

	fmt.Printf("%s = %s\n", key, value)
	return nil
}

func runConfigList(cmd *cobra.Command, args []string) error {
	fmt.Println("Current configuration:")
	fmt.Println()

	settings := map[string]string{
		"api-key":          viper.GetString("api-key"),
		"whisper-model":    viper.GetString("whisper-model"),
		"default-duration": viper.GetString("default-duration"),
		"default-quality":  viper.GetString("default-quality"),
		"temp-dir":         viper.GetString("temp-dir"),
	}

	for key, value := range settings {
		displayValue := value
		if value == "" {
			displayValue = "(not set)"
		} else if key == "api-key" && len(value) > 8 {
			displayValue = value[:4] + "..." + value[len(value)-4:]
		}
		
		fmt.Printf("  %-16s = %s\n", key, displayValue)
	}

	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("\nConfig file: %s\n", configFile)
	}

	return nil
}

func runConfigReset(cmd *cobra.Command, args []string) error {
	fmt.Print("⚠️  This will reset all configuration to defaults. Continue? (y/N): ")
	
	var response string
	fmt.Scanln(&response)
	
	if response != "y" && response != "Y" {
		fmt.Println("Configuration reset cancelled.")
		return nil
	}

	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		if err := os.Remove(configFile); err != nil {
			return fmt.Errorf("failed to remove config file: %w", err)
		}
	}

	fmt.Println("✅ Configuration reset to defaults.")
	return nil
}

func validateConfigKey(key string) error {
	validKeys := []string{
		"api-key",
		"whisper-model", 
		"default-duration",
		"default-quality",
		"temp-dir",
	}

	for _, validKey := range validKeys {
		if key == validKey {
			return nil
		}
	}

	return fmt.Errorf("invalid configuration key: %s (valid keys: %v)", key, validKeys)
}