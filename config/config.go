package config

var (
	// Config is a struct that provides all covey configuration details.
	Config = struct {
		// revive:disable:line-length-limit
		Daemon struct {
			Host          string `default:"" flag:"host||Resolvable address for this machine (If empty covey will listen on every address)." env:"COVEY_HOST"`
			Port          string `default:"8080" flag:"port||Port to expose the covey daemon on. (default '8080')" env:"COVEY_PORT"`
			PluginsFolder string `default:"/usr/lib64/covey/plugins" flag:"plugins-folder||Folder where plugins are located. (default '/usr/lib64/covey/plugins')" env:"COVEY_PLUGINS_FOLDER"`
		}
		DB struct {
			Username string `default:"postgres" flag:"postgres-username||The username used to login to the postgres database. (default 'postgres')" env:"COVEY_POSTGRES_USERNAME"`
			Password string `default:"" flag:"postgres-password||The password used to login to the postgres database." env:"COVEY_POSTGRES_PASSWORD"`
			Host     string `default:"127.0.0.1" flag:"postgres-host||The Postgres host (default '127.0.0.1')" env:"COVEY_POSTGRES_HOST"`
			Port     string `default:"5432" flag:"postgres-port||The Postgres port (default '5432')" env:"COVEY_POSTGRES_PORT"`
			Database string `default:"covey" flag:"postgres-database||The database (default 'covey')" env:"COVEY_POSTGRES_DATABASE"`
		}
		// revive:enable:line-length-limit
	}{}
)

// InitConfig initializes the configuration values.
func InitConfig() error {
	configurator, err := configuration.New(&Config, []configuration.Provider{
		configuration.NewFlagProvider(&Config),
		configuration.NewEnvProvider(),
		configuration.NewFileProvider("/etc/covey/config.yml"),
		configuration.NewDefaultProvider(),
	}, false, false)
	if err != nil {
		return err
	}

	configurator.InitValues()
	return nil
}
