package main

import (
	"fmt"
	"os"
	"slices"
)

type CodeInitArgs struct {
	BaseArgs
	Files []string `arg:"-c,--compose-file,separate,required" help:"Compose files to load"`
	Host  string   `arg:"-d,--host,required" help:"Domain used for routing"`
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
	deployProject := buildDeployProject()
	for name, svc := range project.Services {
		deployProject.Services[name] = createDeployService(svc, slices.Contains(imageNames, svc.Image))
	}

	mpFile := Config{
		Project: ProjectConfig{
			Files: args.Files,
		},
		Api: ApiConfig{
			Host:  fmt.Sprintf("tcp://%s:2376", args.Host),
			Certs: ".tls",
		},
		Deploy: DeployConfig{
			Version: "0.1.0",
		},
	}

	err = os.MkdirAll(args.Env, 0755)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to create directory: %s", args.Env))

	mpPath, err := saveConfig(args.Env, mpFile)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", mpPath))
	fmt.Println("✅ ", mpPath)

	deployPath, err := saveProject(args.Env, deployProject)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", deployPath))
	fmt.Println("✅ ", deployPath)

	serverPath, err := caddyCreateMiniPaasServerFile(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", serverPath))
	fmt.Println("✅ ", serverPath)

}
