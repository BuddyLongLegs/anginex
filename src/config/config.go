package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

func GetConfig() Config {
	configFile := flag.String("f", "", "config file")
	flag.Parse()

	// if config file is provided via flag
	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		// config file: anginex.yml
		viper.SetConfigName("anginex")
		viper.SetConfigType("yaml")

		// supported config paths
		viper.AddConfigPath("/etc/anginex/")
		viper.AddConfigPath("$HOME/.anginex")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var config Config
	err = viper.Unmarshal(&config)

	if err != nil {
		panic(fmt.Errorf("fatal error unmarshal config file: %w", err))
	}

	return config
}
