package main

import (
	"encoding/base64"
	"github.com/goccy/go-yaml"
	"strings"
)

type MiniPaaSConfig struct {
	PostgresHost string `yaml:"postgres_host"`
	PostgresUser string `yaml:"postgres_user"`
	PostgresDB   string `yaml:"postgres_db"`
}

func loadMiniPaaSConfig(verbose bool) (MiniPaaSConfig, error) {
	output, err := runCommandOutput([]string{"docker", "config", "inspect", "minipaas__config", "--format", "{{json .Spec.Data}}"}, verbose)
	if err != nil {
		return MiniPaaSConfig{}, err
	}

	encodedData := strings.TrimSpace(output)
	encodedData = encodedData[1 : len(encodedData)-1]
	println(encodedData)
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return MiniPaaSConfig{}, err
	}

	var config MiniPaaSConfig
	if err = yaml.Unmarshal(decodedData, &config); err != nil {
		return MiniPaaSConfig{}, err
	}

	return config, nil
}
