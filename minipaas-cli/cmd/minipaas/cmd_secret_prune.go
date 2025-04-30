package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type SecretPruneArgs struct {
	BaseArgs
	Delete bool `arg:"--delete" help:"Actually delete the unused secrets; if not set, only list them."`
}

func (args *SecretPruneArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Error loading configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	output, err := runCommandOutput([]string{"docker", "secret", "ls", "--format", "{{.Name}}"}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Error listing secrets"))
	allSecrets := strings.Fields(output)
	if len(allSecrets) == 0 {
		fmt.Println("⚠️ No secrets found.")
		return
	}

	var inspectOutput string
	usedSecrets := make(map[string]bool)
	servicesOutput, err := runCommandOutput([]string{"docker", "service", "ls", "--format", "{{.ID}}"}, args.Verbose)
	checkErrorPanic(err, "❌ Error listing services")
	serviceIDs := strings.Fields(servicesOutput)
	for _, serviceID := range serviceIDs {
		inspectOutput, err = runCommandOutput([]string{"docker", "service", "inspect", serviceID}, args.Verbose)
		if err != nil {
			log.Printf(fmt.Sprintf("⚠️ Error inspecting service %s: %v", serviceID, err))
			continue
		}

		var services []DockerService
		if err = json.Unmarshal([]byte(inspectOutput), &services); err != nil {
			log.Printf(fmt.Sprintf("⚠️ Error parsing JSON for service %s: %v", serviceID, err))
			continue
		}
		if len(services) == 0 {
			continue
		}
		for _, secretEntry := range services[0].Spec.TaskTemplate.ContainerSpec.Secrets {
			usedSecrets[secretEntry.SecretName] = true
		}
	}

	var unusedSecrets []string
	for _, s := range allSecrets {
		if !usedSecrets[s] {
			unusedSecrets = append(unusedSecrets, s)
		}
	}

	fmt.Println("The following secrets are unused:")
	for _, secret := range unusedSecrets {
		fmt.Println(secret)
	}

	if args.Delete {
		fmt.Println("Deleting unused secrets...")
		for _, secret := range unusedSecrets {
			if err = configInternalDelete(secret, args.Verbose); err != nil {
				log.Printf(fmt.Sprintf("❌ Failed to delete secret %s: %v", secret, err))
			} else {
				fmt.Printf("✅ Deleted secret: %s\n", secret)
				_, err = updateDeployFileRemoveSecret(args.Env, secret)
				if err != nil {
					log.Printf("❌ Failed to update deploy file for config %s: %v", secret, err)
				}
			}
		}
	} else if len(unusedSecrets) > 0 {
		fmt.Println("Run with --delete to actually remove the above configs.")
	}
}
