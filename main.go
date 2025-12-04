package main

import (
	"errors"
	"log"

	"github.com/spf13/viper"
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/kunwardeep/paralleltest/pkg/paralleltest"
)

func main() {
	viper.SetConfigName(".paralleltest")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.paralleltest")
	viper.AddConfigPath(".")

	// Read in config, ignore if the file isn't found and use defaults.
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Fatalf("failed to parse config: %v", err)
		}
	}

	var cfg paralleltest.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	singlechecker.Main(paralleltest.NewAnalyzer(cfg))
}
