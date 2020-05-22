package config

import (
	"fmt"
	"os"

	"github.com/fulldump/goconfig"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// GetConfig creates Config struct and fills it fields
func GetConfig(appInfo AppInfo, secrets Secrets) (*Config, error) {
	configName := GetEnvOrDefault("CONFIG_NAME", "dev")

	file := fmt.Sprintf("./config/%s.yml", configName)

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to read config file %s. %s", file, err.Error())
	}

	conf := getDefault(appInfo)
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal config file %s. %s", file, err.Error())
	}

	goconfig.Read(conf)

	// decrypt := func(crypted string) (string, error) {
	// 	decrypted, err := secrets.Decrypt(crypted)
	// 	if err != nil {
	// 		return "", fmt.Errorf("Unable to decrypt secret: %s", err.Error())
	// 	}

	// 	return decrypted, nil
	// }

	return conf, nil
}

// GetEnvOrDefault reads environment variable and returns value
// if there is no environment variable present then def value is returned
func GetEnvOrDefault(key string, def string) string {
	fromEnv := os.Getenv(key)
	if len(fromEnv) == 0 {
		return def
	}
	return fromEnv
}

func getDefault(appInfo AppInfo) *Config {
	return &Config{
		Port: 3000,
		Log: Log{
			Level: "INFO",
		},
		Cors: CorsConfig{
			AllowedOrigins: []string{"*"},
		},
		AppInfo: appInfo,
	}
}
