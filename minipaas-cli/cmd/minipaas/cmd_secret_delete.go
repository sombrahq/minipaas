package main

import (
	"fmt"
	"path/filepath"

	"github.com/compose-spec/compose-go/v2/types"
)

type SecretDeleteArgs struct {
	BaseArgs
	Name string `arg:"positional,required" help:"The name of the Docker secret to delete."`
}

func (args *SecretDeleteArgs) Run() {
	configFile := filepath.Join(args.Env, "minipaas.yaml")
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Error loading configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	err = secretInternalDelete(args.Name, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Error deleting secret: %s", args.Name))
	fmt.Printf("✅ Secret deleted: %s\n", args.Name)

	fn, err := updateDeployFileRemoveSecret(args.Env, args.Name)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to update deploy file: %s", fn))
	fmt.Printf("✅ Updated deploy file: %s\n", fn)
}

func secretInternalDelete(name string, verbose bool) error {
	return runCommand([]string{"docker", "secret", "rm", name}, verbose)
}

func removeSecretFromDeployProject(project *types.Project, secretName string) bool {
	changed := false

	if project.Secrets != nil {
		if _, exists := project.Secrets[secretName]; exists {
			delete(project.Secrets, secretName)
			changed = true
			fmt.Printf("✅ Removed secret %s from top-level secrets\n", secretName)
		}
	}

	for svcName, svc := range project.Services {
		if len(svc.Secrets) > 0 {
			var filteredSecrets []types.ServiceSecretConfig
			for _, sec := range svc.Secrets {
				if sec.Source != secretName {
					filteredSecrets = append(filteredSecrets, sec)
				}
			}
			if len(filteredSecrets) != len(svc.Secrets) {
				svc.Secrets = filteredSecrets
				project.Services[svcName] = svc
				changed = true
				fmt.Printf("✅ Removed secret %s from service %s\n", secretName, svc.Name)
			}
		}
	}
	return changed
}

func updateDeployFileRemoveSecret(env, secretName string) (string, error) {
	project, fn, err := loadProject(env)
	if err != nil {
		return fn, err
	}

	if project.Secrets == nil && len(project.Secrets) == 0 {
		// Nothing to update.
		return fn, nil
	}

	if !removeSecretFromDeployProject(project, secretName) {
		return fn, nil
	}

	return saveProject(env, project)
}
