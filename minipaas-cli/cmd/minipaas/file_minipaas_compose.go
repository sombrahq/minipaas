package main

import (
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const deployFile = "compose.minipaas.yaml"

func createDeployService(svc types.ServiceConfig, requiresBuild bool) types.ServiceConfig {
	srv := types.ServiceConfig{
		Networks: map[string]*types.ServiceNetworkConfig{
			"minipaas_network": {},
		},
	}
	if requiresBuild {
		srv.Image = buildCommonImage(svc.Image)
	}
	return srv
}

func composeEnsureDeploy(project *types.Project, serviceName string) {
	svc, ok := project.Services[serviceName]
	if !ok {
		svc = createDeployService(project.Services[serviceName], false)
		project.Services[serviceName] = svc
	}
}

func loadProject(env string) (*types.Project, string, error) {
	ctx := context.Background()
	fn := filepath.Join(env, deployFile)

	options, err := cli.NewProjectOptions(
		[]string{fn},
		// dirty trick. MINIPAAS_DEPLOY_VERSION needs to be defined
		cli.WithEnv([]string{"MINIPAAS_DEPLOY_VERSION=${MINIPAAS_DEPLOY_VERSION}"}),
		cli.WithName("minipaas"),
		cli.WithResolvedPaths(false),
		cli.WithConsistency(false),
		cli.WithoutEnvironmentResolution,
	)
	if err != nil {
		return nil, env, err
	}

	project, err := options.LoadProject(ctx)
	if err != nil {
		return nil, env, err
	}

	if project.Networks != nil {
		if _, ok := project.Networks["default"]; ok {
			delete(project.Networks, "default")
		}
	}

	for name, svc := range project.Services {
		if svc.Networks != nil {
			if _, exists := svc.Networks["default"]; exists {
				delete(svc.Networks, "default")
				project.Services[name] = svc
			}
		}
	}

	return project, env, err
}

func saveProject(env string, project *types.Project) (string, error) {
	fn := filepath.Join(env, deployFile)

	project.Name = ""
	data, err := project.MarshalYAML()
	if err != nil {
		return fn, err
	}
	return fn, os.WriteFile(fn, data, 0644)
}

func buildCommonImage(image string) string {
	version := "${MINIPAAS_DEPLOY_VERSION}"

	if idx := strings.LastIndex(image, ":"); idx != -1 {
		image = image[:idx]
	}

	if !strings.Contains(image, "/") {
		image = "registry:5000/" + image
	}

	return image + ":" + version
}

func buildDeployProject() *types.Project {
	return &types.Project{
		Services: make(types.Services),
		Networks: types.Networks{
			"minipaas_network": types.NetworkConfig{
				External: true,
			},
		},
	}
}

func addComposeResilientDeploy(project *types.Project, service string, port string) error {
	svc, ok := project.Services[service]
	if !ok {
		return fmt.Errorf("service %s not found in project", service)
	}

	replicas := 2

	healthInterval := types.Duration(10 * time.Second)
	healthTimeout := types.Duration(10 * time.Second)
	healthReties := uint64(5)
	healthStart := types.Duration(10 * time.Second)

	updateParallelism := uint64(1)
	updateDelay := types.Duration(10 * time.Second)

	rollbackParallelism := uint64(0)

	restartDelay := types.Duration(10 * time.Second)
	rollbackDelay := types.Duration(10 * time.Second)

	svc.HealthCheck = &types.HealthCheckConfig{
		Test:        []string{"CMD-SHELL", fmt.Sprintf("wget -qO- --spider http://127.0.0.1:%s", port)},
		Interval:    &healthInterval,
		Timeout:     &healthTimeout,
		Retries:     &healthReties,
		StartPeriod: &healthStart,
	}

	svc.Deploy = &types.DeployConfig{
		Mode:     "replicated",
		Replicas: &replicas,
		UpdateConfig: &types.UpdateConfig{
			Parallelism:   &updateParallelism,
			Order:         "start-first",
			FailureAction: "rollback",
			Delay:         updateDelay,
		},
		RestartPolicy: &types.RestartPolicy{
			Condition: "any",
			Delay:     &restartDelay,
		},
		RollbackConfig: &types.UpdateConfig{
			Parallelism: &rollbackParallelism,
			Order:       "stop-first",
			Delay:       rollbackDelay,
		},
	}

	project.Services[service] = svc

	return nil
}

func addComposeWorkerDeploy(project *types.Project, service string) error {
	svc, ok := project.Services[service]
	if !ok {
		return fmt.Errorf("service %s not found in project", service)
	}

	replicas := 1

	updateParallelism := uint64(1)
	updateDelay := types.Duration(10 * time.Second)

	rollbackParallelism := uint64(0)

	restartDelay := types.Duration(10 * time.Second)
	rollbackDelay := types.Duration(10 * time.Second)

	svc.Deploy = &types.DeployConfig{
		Mode:     "replicated",
		Replicas: &replicas,
		UpdateConfig: &types.UpdateConfig{
			Parallelism:   &updateParallelism,
			Order:         "start-first",
			FailureAction: "rollback",
			Delay:         updateDelay,
		},
		RestartPolicy: &types.RestartPolicy{
			Condition: "any",
			Delay:     &restartDelay,
		},
		RollbackConfig: &types.UpdateConfig{
			Parallelism: &rollbackParallelism,
			Order:       "stop-first",
			Delay:       rollbackDelay,
		},
	}

	project.Services[service] = svc

	return nil
}

func addComposeJobDeploy(project *types.Project, serviceName string) error {
	svc, ok := project.Services[serviceName]
	if !ok {
		return fmt.Errorf("service %s not found in project", serviceName)
	}

	replicas := 1
	updateParallelism := uint64(0)
	updateDelay := types.Duration(10 * time.Second)

	restartDelay := types.Duration(10 * time.Second)
	restartAttempts := uint64(10)

	svc.Deploy = &types.DeployConfig{
		Mode:     "replicated",
		Replicas: &replicas,
		UpdateConfig: &types.UpdateConfig{
			Parallelism:   &updateParallelism,
			Order:         "stop-first",
			FailureAction: "pause",
			Delay:         updateDelay,
		},
		RestartPolicy: &types.RestartPolicy{
			Condition:   "on-failure",
			Delay:       &restartDelay,
			MaxAttempts: &restartAttempts,
		},
	}

	project.Services[serviceName] = svc
	return nil
}

func addComposeCronDeploy(project *types.Project, serviceName string, cron string) error {
	// https://crazymax.dev/swarm-cronjob/
	svc, ok := project.Services[serviceName]
	if !ok {
		return fmt.Errorf("service %s not found in project", serviceName)
	}

	replicas := 0

	svc.Deploy = &types.DeployConfig{
		Mode:     "replicated",
		Replicas: &replicas,
		RestartPolicy: &types.RestartPolicy{
			Condition: "none",
		},
		Labels: map[string]string{
			"swarm.cronjob.enable":       "true",
			"swarm.cronjob.schedule":     cron,
			"swarm.cronjob.skip-running": "true",
		},
	}

	project.Services[serviceName] = svc
	return nil
}

func addComposeConfig(project *types.Project, config, name string, services []string) error {
	if project.Configs == nil {
		project.Configs = make(map[string]types.ConfigObjConfig)
	}

	project.Configs[config] = types.ConfigObjConfig{
		External: true,
	}

	for _, svcName := range services {
		svc, ok := project.Services[svcName]
		if !ok {
			return fmt.Errorf("service %s not found in project", svcName)
		}
		if svc.Configs == nil {
			svc.Configs = []types.ServiceConfigObjConfig{}
		}

		exists := false
		for _, sc := range svc.Configs {
			if sc.Source == config {
				exists = true
				break
			}
		}
		if !exists {
			svc.Configs = append(svc.Configs, types.ServiceConfigObjConfig{
				Source: config,
				Target: name,
			})
		}
		project.Services[svcName] = svc
	}

	return nil
}

func addComposeSecret(project *types.Project, secret, name string, services []string) error {
	if project.Secrets == nil {
		project.Secrets = make(map[string]types.SecretConfig)
	}

	project.Secrets[secret] = types.SecretConfig{
		External: true,
	}

	for _, svcName := range services {
		svc, ok := project.Services[svcName]
		if !ok {
			return fmt.Errorf("service %s not found in project", svcName)
		}
		if svc.Secrets == nil {
			svc.Secrets = []types.ServiceSecretConfig{}
		}

		exists := false
		for _, sc := range svc.Secrets {
			if sc.Source == secret {
				exists = true
				break
			}
		}
		if !exists {
			svc.Secrets = append(svc.Secrets, types.ServiceSecretConfig{
				Source: secret,
				Target: name,
			})
		}
		project.Services[svcName] = svc
	}

	return nil
}

func updateDeployFileRemoveConfig(env, configName string) (string, error) {
	project, fn, err := loadProject(env)
	if err != nil {
		return fn, err
	}

	if project.Configs == nil && len(project.Services) == 0 {
		return fn, nil
	}

	if !removeConfigFromDeployProject(project, configName) {
		return fn, nil
	}

	fn, err = saveProject(env, project)
	return fn, err
}

func removeConfigFromDeployProject(project *types.Project, configName string) bool {
	changed := false

	if project.Configs != nil {
		if _, exists := project.Configs[configName]; exists {
			delete(project.Configs, configName)
			changed = true
			fmt.Printf("✅ Removed config %s from top-level configs\n", configName)
		}
	}

	for svcName, svc := range project.Services {
		if len(svc.Configs) > 0 {
			var filteredConfigs []types.ServiceConfigObjConfig
			for _, cfg := range svc.Configs {
				if cfg.Source != configName {
					filteredConfigs = append(filteredConfigs, cfg)
				}
			}
			if len(filteredConfigs) != len(svc.Configs) {
				svc.Configs = filteredConfigs
				project.Services[svcName] = svc
				changed = true
				fmt.Printf("✅ Removed config %s from service %s\n", configName, svc.Name)
			}
		}
	}
	return changed
}
