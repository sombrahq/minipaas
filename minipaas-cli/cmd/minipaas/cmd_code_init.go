package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type CodeInitArgs struct {
	BaseArgs
	Files  []string `arg:"-c,--compose-file,separate,required" help:"Compose files to load"`
	Domain string   `arg:"-d,--domain" help:"Domain used for routing"`
	Local  bool     `arg:"--local" help:"Use local Docker daemon" default:"true"`
}

func (args *CodeInitArgs) Run() {
	project, err := composeLoadProject(args.Files)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load project files: %s", args.Files))

	imageNames := make([]string, 0)
	for _, srv := range project.Services {
		if srv.Build != nil {
			imageNames = append(imageNames, srv.Image)
		}
	}
	appsCompose := buildDeployProject()
	for name, svc := range project.Services {
		appsCompose.Services[name] = createDeployService(svc, slices.Contains(imageNames, svc.Image))
	}

	var api ApiConfig
	if args.Domain != "" {
		args.Local = true
	}
	if args.Local {
		api = ApiConfig{
			Local: true,
		}
	} else {
		api = ApiConfig{
			Host:  fmt.Sprintf("tcp://%s:2376", args.Domain),
			Certs: ".tls",
			Local: false,
		}
	}

	err = os.MkdirAll(args.Env, 0755)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to create directory: %s", args.Env))

	swarmAppsPath, err := saveProject(args.Env, appsCompose)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", swarmAppsPath))
	fmt.Println("✅ ", swarmAppsPath)

	swarmCommonFile := filepath.Join(args.Env, "compose.common.yml")
	swarmCommonContent, err := readEmbeddedSwarm("common.yml")
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load embed file: %s", "common.yml"))
	err = os.WriteFile(swarmCommonFile, swarmCommonContent, 0644)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", swarmCommonFile))
	fmt.Println("✅ ", swarmCommonFile)

	swarmRegistryFile := filepath.Join(args.Env, "compose.registry.yml")
	swarmRegistryContent, err := readEmbeddedSwarm("registry.yml")
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load embed file: %s", "registry.yml"))
	err = os.WriteFile(swarmRegistryFile, swarmRegistryContent, 0644)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", swarmRegistryFile))
	fmt.Println("✅ ", swarmRegistryFile)

	// optional Postgres service override file
	swarmPostgresFile := filepath.Join(args.Env, "compose.postgres.yml")
	swarmPostgresContent, err := readEmbeddedSwarm("postgres.yaml")
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load embed file: %s", "postgres.yaml"))
	err = os.WriteFile(swarmPostgresFile, swarmPostgresContent, 0644)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", swarmPostgresFile))
	fmt.Println("✅ ", swarmPostgresFile)

	swarmCaddyFile := filepath.Join(args.Env, "compose.caddy.yml")
	swarmCaddyContent, err := readEmbeddedSwarm("caddy.yml")
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load embed file: %s", "caddy.yml"))
	err = os.WriteFile(swarmCaddyFile, swarmCaddyContent, 0644)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", swarmCaddyFile))
	fmt.Println("✅ ", swarmCaddyFile)

	caddyConfFile := filepath.Join(args.Env, "caddy.json")
	caddyConfContent, err := readEmbeddedSwarm("caddy.json")
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load embed file: %s", "caddy.json"))
	err = os.WriteFile(caddyConfFile, caddyConfContent, 0644)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", caddyConfFile))
	fmt.Println("✅ ", caddyConfFile)

	files := append(args.Files, swarmCommonFile, swarmAppsPath, swarmRegistryFile, swarmPostgresFile, swarmCaddyFile)

	mpFile := Config{
		Project: ProjectConfig{
			Files: files,
		},
		Api: api,
		Deploy: DeployConfig{
			Version: "0.1.0",
		},
	}
	mpPath, err := saveConfig(args.Env, mpFile)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", mpPath))
	fmt.Println("✅ ", mpPath)

}
