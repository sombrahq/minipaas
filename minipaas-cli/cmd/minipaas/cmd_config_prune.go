package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type ConfigPruneArgs struct {
	BaseArgs
	Delete bool `arg:"--delete" help:"Actually delete the unused configs; if not set, only list them."`
}

func (args *ConfigPruneArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Error loading configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	output, err := runCommandOutput([]string{"docker", "config", "ls", "--format", "{{.Name}}"}, args.Verbose)
	checkErrorPanic(err, "❌ Error listing configs")
	allConfigs := strings.Fields(output)
	if len(allConfigs) == 0 {
		fmt.Println("No configs found.")
		return
	}

	var inspectOutput string
	usedConfigs := make(map[string]bool)
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
		for _, configEntry := range services[0].Spec.TaskTemplate.ContainerSpec.Configs {
			usedConfigs[configEntry.ConfigName] = true
		}
	}

	var unusedConfigs []string
	for _, config := range allConfigs {
		if !usedConfigs[config] && config != "minipaas__config" {
			unusedConfigs = append(unusedConfigs, config)
		}
	}

	fmt.Println("The following configs are unused:")
	for _, config := range unusedConfigs {
		fmt.Println(config)
	}

	if args.Delete {
		fmt.Println("Deleting unused configs...")
		for _, config := range unusedConfigs {
			if err = configInternalDelete(config, args.Verbose); err != nil {
				log.Printf("❌ Failed to delete config %s: %v", config, err)
			} else {
				fmt.Printf("✅ Deleted config: %s\n", config)
				_, err = updateDeployFileRemoveConfig(args.Env, config)
				if err != nil {
					log.Printf("❌ Failed to update deploy file for config %s: %v", config, err)
				}
			}
		}
	} else if len(unusedConfigs) > 0 {
		fmt.Println("Run with --delete to actually remove the above configs.")
	}
}
