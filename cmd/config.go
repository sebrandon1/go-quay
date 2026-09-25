package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v3"
)

var userConfigDir = os.UserConfigDir

// appConfig holds values loaded from the config file.
// Fields are applied as defaults — CLI flags and environment variables always take priority.
type appConfig struct {
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	QuayURL   string `yaml:"quay-url"`
}

// appCfg is initialized at package load time, before any init() functions run.
// This ensures config values are available when flags are registered.
var appCfg = loadConfig()

var (
	configNamespace          string
	configInitNonInteractive bool
	configShowToken          bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View and initialize CLI configuration",
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the configuration file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := configFilePath()
		if path == "" {
			return fmt.Errorf("determining config file path")
		}
		_, err := fmt.Fprintln(cmd.OutOrStdout(), path)
		return err
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current configuration",
	Long:  `Show the saved configuration. The token is redacted unless --show-token is supplied.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := readConfig(configFilePath())
		if err != nil {
			return fmt.Errorf("reading config: %w", err)
		}
		if !configShowToken && cfg.Token != "" {
			cfg.Token = "********"
		}
		output, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("marshaling config: %w", err)
		}
		_, err = cmd.OutOrStdout().Write(output)
		return err
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the CLI configuration file",
	Long: `Create or update the CLI configuration file.

Use --token and --namespace for non-interactive setup; --quay-url is optional. When stdin is not a terminal, --token and --namespace are required.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := configFilePath()
		if path == "" {
			return fmt.Errorf("determining config file path")
		}
		cfg, err := configInitDefaults(cmd)
		if err != nil {
			return err
		}

		interactive := !configInitNonInteractive && isTerminal(cmd.InOrStdin())
		if !interactive {
			if !flagChanged(cmd, "token") || cfg.Token == "" {
				return fmt.Errorf("non-interactive config init requires --token")
			}
			if !cmd.Flags().Changed("namespace") || cfg.Namespace == "" {
				return fmt.Errorf("non-interactive config init requires --namespace")
			}
		} else {
			cfg, err = promptForConfig(cmd.InOrStdin(), cmd.ErrOrStderr(), cfg)
			if err != nil {
				return err
			}
		}

		if cfg.Token == "" {
			return fmt.Errorf("token is required")
		}
		if cfg.Namespace == "" {
			return fmt.Errorf("namespace is required")
		}
		if cfg.QuayURL == "" {
			cfg.QuayURL = lib.DefaultQuayURL
		}
		if err := writeConfig(path, cfg); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}
		appCfg = cfg
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Configuration written to %s\n", path)
		return err
	},
}

// loadConfig reads the config file from the user's config directory.
// Returns an empty config if the file doesn't exist or can't be parsed.
func loadConfig() appConfig {
	cfg, err := readConfig(configFilePath())
	if err != nil {
		return appConfig{}
	}
	return cfg
}

func readConfig(path string) (appConfig, error) {
	if path == "" {
		return appConfig{}, fmt.Errorf("config path is empty")
	}
	data, err := os.ReadFile(path) // #nosec G304 -- path is derived from os.UserConfigDir()
	if errors.Is(err, os.ErrNotExist) {
		return appConfig{}, nil
	}
	if err != nil {
		return appConfig{}, fmt.Errorf("reading %s: %w", path, err)
	}
	var cfg appConfig
	if err := parseConfig(data, &cfg); err != nil {
		return appConfig{}, fmt.Errorf("parsing %s: %w", path, err)
	}
	return cfg, nil
}

// parseConfig unmarshals YAML config data into the given appConfig.
func parseConfig(data []byte, cfg *appConfig) error {
	return yaml.Unmarshal(data, cfg)
}

// configFilePath returns the path to the config file.
// Uses os.UserConfigDir() to determine the platform-appropriate config directory:
//   - Linux:   ~/.config/go-quay/config.yaml
//   - macOS:   ~/Library/Application Support/go-quay/config.yaml
//   - Windows: %AppData%/go-quay/config.yaml
func configFilePath() string {
	dir, err := userConfigDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, cliName, "config.yaml")
}

func configInitDefaults(cmd *cobra.Command) (appConfig, error) {
	cfg, err := readConfig(configFilePath())
	if err != nil {
		return appConfig{}, err
	}
	cfg.Token = firstNonEmpty(os.Getenv("QUAY_TOKEN"), cfg.Token)
	cfg.QuayURL = firstNonEmpty(os.Getenv("QUAY_URL"), cfg.QuayURL, lib.DefaultQuayURL)
	if flagChanged(cmd, "token") {
		cfg.Token = cmd.Flag("token").Value.String()
	}
	if cmd.Flags().Changed("namespace") {
		cfg.Namespace = configNamespace
	}
	if flagChanged(cmd, "quay-url") {
		cfg.QuayURL = cmd.Flag("quay-url").Value.String()
	}
	return cfg, nil
}

func promptForConfig(reader io.Reader, writer io.Writer, cfg appConfig) (appConfig, error) {
	input := bufio.NewReader(reader)
	tokenPrompt := "Quay token: "
	if cfg.Token != "" {
		tokenPrompt = "Quay token (leave blank to keep current value): "
	}
	token, err := readConfigPrompt(input, writer, tokenPrompt, cfg.Token)
	if err != nil {
		return appConfig{}, fmt.Errorf("reading token: %w", err)
	}
	if token != "" {
		cfg.Token = token
	}

	namespacePrompt := "Namespace: "
	if cfg.Namespace != "" {
		namespacePrompt = fmt.Sprintf("Namespace [%s]: ", cfg.Namespace)
	}
	namespace, err := readConfigPrompt(input, writer, namespacePrompt, cfg.Namespace)
	if err != nil {
		return appConfig{}, fmt.Errorf("reading namespace: %w", err)
	}
	if namespace != "" {
		cfg.Namespace = namespace
	}

	urlPrompt := fmt.Sprintf("Quay API URL [%s]: ", cfg.QuayURL)
	quayURL, err := readConfigPrompt(input, writer, urlPrompt, cfg.QuayURL)
	if err != nil {
		return appConfig{}, fmt.Errorf("reading Quay URL: %w", err)
	}
	if quayURL != "" {
		cfg.QuayURL = quayURL
	}
	return cfg, nil
}

func readConfigPrompt(reader *bufio.Reader, writer io.Writer, prompt, defaultValue string) (string, error) {
	if _, err := fmt.Fprint(writer, prompt); err != nil {
		return "", err
	}
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	value := strings.TrimSpace(line)
	if value == "" {
		return defaultValue, nil
	}
	return value, nil
}

func isTerminal(reader io.Reader) bool {
	file, ok := reader.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func writeConfig(path string, cfg appConfig) error {
	if path == "" {
		return fmt.Errorf("config path is empty")
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		return fmt.Errorf("restricting config directory permissions: %w", err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}
	if err := file.Chmod(0o600); err != nil {
		_ = file.Close()
		return fmt.Errorf("restricting config file permissions: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return fmt.Errorf("writing config file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("closing config file: %w", err)
	}
	return nil
}

func init() {
	configCmd.AddCommand(configInitCmd, configPathCmd, configShowCmd)
	configInitCmd.Flags().StringVar(&configNamespace, "namespace", appCfg.Namespace, "Default namespace to save")
	configInitCmd.Flags().BoolVar(&configInitNonInteractive, "non-interactive", false, "Require flags instead of prompting")
	configShowCmd.Flags().BoolVar(&configShowToken, "show-token", false, "Display the saved token without redaction")
}
