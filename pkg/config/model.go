package config

// Config contains configuration data for modules in this project
type Config struct {
	Port      int        `yaml:"port" usage:"Application port"`
	AwsRegion string     `yaml:"awsRegion"`
	Cors      CorsConfig `yaml:"cors"`
	Log       Log
	AppInfo   AppInfo // initialized by appCode not yaml config
}

type Log struct {
	Level string `yaml:"logLevel"`
}

type CorsConfig struct {
	AllowedOrigins []string `yaml:"allowedOrigins"`
}

type AppInfo struct {
	Version   string
	GitCommit string
	Name      string
}

type Secrets interface {
	Decrypt(encodedSecret string) (string, error)
}

func (c Config) LogLevel() string {
	return c.Log.Level
}

func (c Config) Version() string {
	return c.AppInfo.Version
}

func (c Config) GitCommit() string {
	return c.AppInfo.GitCommit
}

func (c Config) Name() string {
	return c.AppInfo.Name
}
