package config

import (
	"github.com/BoRuDar/configuration/v2"
)

var (
	Config = struct {
		Daemon struct {
			Host string `default:"" flag:"host||Resolvable address for this machine." env:"COVEY_HOST"`
			Port int    `default:"8080" flag:"port||Port to expose the covey daemon on. (default "8080")" env:"COVEY_PORT"`
		}
		DB struct {
			Username string `default:"postgres" flag:"postgres-username||The username used to login to the postgres database. (default "postgres")" env:"COVEY_POSTGRES_USERNAME"`
			Password string `default:"" flag:"postgres-password||The password used to login to the postgres database." env:"COVEY_POSTGRES_PASSWORD"`
			Host     string `default:"127.0.0.1" flag:"postgres-host||The Postgres host (default "127.0.0.1")" env:"COVEY_POSTGRES_HOST"`
			Port     int    `default:"5432" flag:"postgres-port||The Postgres port (default "5432")" env:"COVEY_POSTGRES_PORT"`
			Database string `default:"covey" flag:"postgres-database||The database (default "covey") env:"COVEY_POSTGRES_DATABASE"`
		}
	}{}
)

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
