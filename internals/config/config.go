package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`

	App struct {
		Port       int    `yaml:"port" env-required:"true"`
		ServerAddr string `yaml:"server_addr" env-required:"true"`
	} `yaml:"app"`

	Database struct {
		HostLocal  string `yaml:"hostlocal" env-required:"true"`
		HostDocker string `yaml:"hostdocker" env-required:"true"`
		Host       string `yaml:"host"`
		Port       int    `yaml:"port" env-required:"true"`
		User       string `yaml:"user" env-required:"true"`
		Password   string `yaml:"password" env-required:"true"`
		Name       string `yaml:"name" env-required:"true"`
	} `yaml:"database"`

	Migrations struct {
		Path string `yaml:"path" env-required:"true"`
	} `yaml:"migrations"`

	Storage struct {
		Type string `yaml:"type" env-required:"true"`
	} `yaml:"storage"`
}

type Flags struct {
	ConfigPath string
	RunAppType string
}

func MustLoad() *Config {
	flags := parseFlags()

	if flags.ConfigPath == "" {
		panic("config path is empty (use --config or CONFIG_PATH)")
	}

	if _, err := os.Stat(flags.ConfigPath); os.IsNotExist(err) {
		panic("config file not found: " + flags.ConfigPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(flags.ConfigPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	switch flags.RunAppType {
	case "localhost":
		cfg.Database.Host = cfg.Database.HostLocal
	case "docker":
		cfg.Database.Host = cfg.Database.HostDocker
	default:
		panic(fmt.Sprintf("unknown app-type: %q (expected 'localhost' or 'docker')", flags.RunAppType))
	}

	fmt.Println("run app with type:", flags.RunAppType)

	return &cfg
}

func parseFlags() *Flags {
	var f Flags
	flag.StringVar(&f.ConfigPath, "config", "", "path to config file (or use CONFIG_PATH)")
	flag.StringVar(&f.RunAppType, "app-type", "localhost", "type of application to run (with docker or local)")
	flag.Parse()
	return &f
}
