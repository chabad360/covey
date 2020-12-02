package config

import (
	"github.com/BoRuDar/configuration/v3"
	"github.com/go-playground/validator/v10"
	"log"
)

var (
	// Config is a struct  that provides all covey configuration details.
	Config = struct {
		// revive:disable:line-length-limit
		Daemon struct {
			Host string `default:"0.0.0.0" flag:"host||Resolvable address for this machine. (default '0.0.0.0')" env:"COVEY_HOST" validate:"omitempty,ip_addr"` // TODO: Use hostname by default
			Port string `default:"8080" flag:"port||Port to expose the covey daemon on. (default '8080')" env:"COVEY_PORT" validate:"number,required"`
		}
		DB struct {
			Username string `default:"postgres" flag:"postgres-username||The username used to login to the postgres database. (default 'postgres')" env:"COVEY_POSTGRES_USERNAME" validate:"required"`
			Password string `default:"" flag:"postgres-password||The password used to login to the postgres database." env:"COVEY_POSTGRES_PASSWORD"`
			Host     string `default:"localhost" flag:"postgres-host||The Postgres host (default 'localhost')" env:"COVEY_POSTGRES_HOST" validate:"required,hostname_rfc1123"`
			Port     string `default:"5432" flag:"postgres-port||The Postgres port. (default '5432')" env:"COVEY_POSTGRES_PORT" validate:"required,number"`
			Database string `default:"covey" flag:"postgres-database||The database. (default 'covey')" env:"COVEY_POSTGRES_DATABASE" validate:"required"`
		}
		Plugins struct {
			PluginsFolder      string `default:"/usr/lib64/covey/plugins/" flag:"plugins-dir||Directory where plugins are located. (default '/usr/lib64/covey/plugins')" env:"COVEY_PLUGINS_DIRECTORY" validate:"dir,required"`
			PluginsCacheFolder string `default:"/var/cache/covey/plugins/" flag:"plugins-cache-dir||Directory where extracted plugins are located. (default '/var/cache/covey/plugins')" env:"COVEY_PLUGINS_CACHE_DIRECTORY" validate:"dir,required"`
		}
		// revive:enable:line-length-limit
	}{}
)

// InitConfig initializes the configuration values.
func InitConfig() error {
	fp, err := configuration.NewFileProvider("/etc/covey/config.yml")
	if err != nil {
		log.Println(err)
	}
	configurator, err := configuration.New(&Config,
		configuration.NewFlagProvider(&Config),
		configuration.NewEnvProvider(),
		fp,
		configuration.NewDefaultProvider(),
	)
	if err != nil {
		return err
	}
	configurator.SetOnFailFn(func(err error) {
		log.Println(err)
	})

	configurator.InitValues()
	return validate(Config)
}

func validate(cfg interface{}) error {
	val := validator.New()
	return val.Struct(cfg)
}
