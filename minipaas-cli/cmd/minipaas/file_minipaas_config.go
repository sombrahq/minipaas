package main

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"os"
	"path/filepath"
)

type ProjectConfig struct {
	Files []string `yaml:"files"`
}

type ApiConfig struct {
	Host  string `yaml:"host"`
	Certs string `yaml:"tls"`
}

type DeployConfig struct {
	Version string `yaml:"version"`
}

type Config struct {
	Project ProjectConfig `yaml:"project"`
	Api     ApiConfig     `yaml:"api"`
	Deploy  DeployConfig  `yaml:"deploy"`
}

func loadConfig(env string) (Config, string, error) {
	fn := filepath.Join(env, "minipaas.yaml")
	data, err := os.ReadFile(fn)
	if err != nil {
		return Config{}, fn, err
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fn, err
	}

	return cfg, fn, nil
}

func saveConfig(env string, cfg Config) (string, error) {
	fn := filepath.Join(env, "minipaas.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fn, err
	}

	return fn, os.WriteFile(fn, data, 0644)
}

func setApiEnvVars(env string, cfg Config, verbose bool) {
	certsPath := filepath.Join(env, cfg.Api.Certs)
	if verbose {
		fmt.Printf("🔹 Environment variables set:\n"+
			"   DOCKER_CERT_PATH=%s\n"+
			"   DOCKER_HOST=%s\n"+
			"   DOCKER_TLS_VERIFY=1\n"+
			"   MINIPAAS_DEPLOY_VERSION=%s\n",
			certsPath, cfg.Api.Host, cfg.Deploy.Version)
	}

	os.Setenv("DOCKER_CERT_PATH", certsPath)
	os.Setenv("DOCKER_HOST", cfg.Api.Host)
	os.Setenv("DOCKER_TLS_VERIFY", "1")
	os.Setenv("MINIPAAS_DEPLOY_VERSION", cfg.Deploy.Version)
}
