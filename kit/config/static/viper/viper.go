package viper

import (
	kit_config "akeneo-migrator/kit/config/static"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type viperConfig struct{}

// NewViperConfig fetch configurations.
func NewViperConfig() kit_config.ConfigurationLoader {
	return &viperConfig{}
}

// LoadConfiguration load the setup for the configuration object.
func (vp *viperConfig) LoadConfiguration(context string) error {
	if viper.IsSet(context) {
		return nil
	}
	viperContext := *viper.New()

	cwd, _ := os.Getwd()
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	servicePath := strings.ToLower(os.Getenv("CONFIG_PATH"))

	// Set the file name of the configurations file
	configPath := fmt.Sprintf("%s/configs/settings.%s.json", cwd, env)
	if env == "pipeline" {
		_, compilationPath, _, _ := runtime.Caller(0)
		projectPath := filepath.Join(filepath.Dir(compilationPath), "../../../..")
		configPath = fmt.Sprintf("%s/configs/%s/settings.%s.json", projectPath, servicePath, env)
	}

	viperContext.SetConfigName(filepath.Base(configPath))
	viperContext.SetConfigType("json")
	viperContext.AddConfigPath(filepath.Dir(configPath))

	// Enable VIPER to read Environment Variables
	viperContext.AutomaticEnv()
	viperContext.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viperContext.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("file has not been found in the current directory")
		} else {
			return fmt.Errorf("fatal error config file: %s \n", err)
		}
	}
	viper.Set(context, viperContext)

	return nil
}
