package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type ProjectConfig struct {
	Files []string `yaml:"files"`
}

type ApiConfig struct {
	Host  string `yaml:"host,omitempty"`
	Certs string `yaml:"tls,omitempty"`
	Local bool   `yaml:"local,omitempty"`
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
	// Local mode: use default Docker socket, no TLS envs
	if verbose {
		fmt.Printf("🔹 Environment:\n"+
			"   MINIPAAS_DEPLOY_VERSION=%s\n", cfg.Deploy.Version)
	}
	os.Setenv("MINIPAAS_DEPLOY_VERSION", cfg.Deploy.Version)

	if cfg.Api.Local {
		return
	}

	// TLS mode (default for remote Docker APIs)
	certsPath := filepath.Join(env, cfg.Api.Certs)
	os.Setenv("DOCKER_CERT_PATH", certsPath)
	os.Setenv("DOCKER_HOST", cfg.Api.Host)
	os.Setenv("DOCKER_TLS_VERIFY", "1")
	os.Setenv("MINIPAAS_DEPLOY_VERSION", cfg.Deploy.Version)

	if verbose {
		fmt.Printf(
			"   DOCKER_CERT_PATH=%s\n"+
				"   DOCKER_HOST=%s\n"+
				"   DOCKER_TLS_VERIFY=1\n"+
				"   MINIPAAS_DEPLOY_VERSION=%s\n",
			certsPath, cfg.Api.Host, cfg.Deploy.Version)
	}

}
