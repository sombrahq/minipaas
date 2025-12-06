package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/compose-spec/compose-go/v2/types"
)

func TestComposeFilesForEnv(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "compose.common.yml")
	f2 := filepath.Join(dir, "compose.registry.yml")
	f3 := filepath.Join(t.TempDir(), "compose.other.yml")

	cfg := Config{Project: ProjectConfig{Files: []string{f1, f2, f3}}}
	got := composeFilesForEnv(dir, cfg)
	if len(got) != 2 || got[0] != f1 || got[1] != f2 {
		t.Fatalf("composeFilesForEnv wrong order/selection: %#v", got)
	}
}

func TestComposeFilesForEnv_RelativePaths(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.yml")
	b := filepath.Join(dir, "sub", "b.yml")
	cfg := Config{Project: ProjectConfig{Files: []string{a, b, "/tmp/outside.yml"}}}
	got := composeFilesForEnv(dir, cfg)
	if len(got) != 2 || got[0] != a || got[1] != b {
		t.Fatalf("unexpected compose files: %#v", got)
	}
}

func TestGroupServicesByComposeFile(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.yml")
	f2 := filepath.Join(dir, "b.yml")

	if err := os.WriteFile(f1, []byte("version: '3.9'\nservices:\n  api:\n    image: x\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte("version: '3.9'\nservices:\n  web:\n    image: y\n"), 0644); err != nil {
		t.Fatal(err)
	}

	files := []string{f1, f2}
	svcMap, missing := groupServicesByComposeFile(files, []string{"api", "web", "worker"})
	if len(missing) != 1 || missing[0] != "worker" {
		t.Fatalf("missing mismatch: %#v", missing)
	}
	if len(svcMap[f1]) != 1 || svcMap[f1][0] != "api" {
		t.Fatalf("f1 map mismatch: %#v", svcMap[f1])
	}
	if len(svcMap[f2]) != 1 || svcMap[f2][0] != "web" {
		t.Fatalf("f2 map mismatch: %#v", svcMap[f2])
	}
}

func TestLoadComposeFile_StripsDefaultNetwork(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "one.yml")
	// include default network references to ensure they are removed
	data := "version: '3.9'\nnetworks:\n  default: {}\nservices:\n  web:\n    image: x\n    networks:\n      - default\n"
	if err := os.WriteFile(f, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	p, _, err := loadComposeFile(f)
	if err != nil {
		t.Fatalf("loadComposeFile: %v", err)
	}
	if _, ok := p.Networks["default"]; ok {
		t.Fatalf("default network should be stripped")
	}
	if svc := p.Services["web"]; svc.Networks != nil {
		if _, ok := svc.Networks["default"]; ok {
			t.Fatalf("service default network should be stripped")
		}
	}
}

func TestSaveComposeFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "round.yml")
	p := &types.Project{Services: make(types.Services)}
	p.Services["a"] = types.ServiceConfig{Image: "busybox"}
	if _, err := saveComposeFile(f, p); err != nil {
		t.Fatalf("saveComposeFile: %v", err)
	}
	if _, err := os.Stat(f); err != nil {
		t.Fatalf("file not written: %v", err)
	}
	// reload
	p2, _, err := loadComposeFile(f)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if _, ok := p2.Services["a"]; !ok {
		t.Fatalf("service missing after roundtrip")
	}
}

func TestLoadAndSaveProject_AppsFile(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("MINIPAAS_DEPLOY_VERSION", "1.2.3")
	defer os.Unsetenv("MINIPAAS_DEPLOY_VERSION")
	// write compose.apps.yaml expected by loadProject/saveProject
	apps := filepath.Join(dir, "compose.apps.yaml")
	content := "version: '3.9'\nnetworks:\n  default: {}\nservices:\n  api:\n    image: repo/app:${MINIPAAS_DEPLOY_VERSION}\n    networks: [default]\n"
	if err := os.WriteFile(apps, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, _, err := loadProject(dir)
	if err != nil {
		t.Fatalf("loadProject: %v", err)
	}
	// default network removed
	if _, ok := p.Networks["default"]; ok {
		t.Fatalf("default network not stripped")
	}
	// save
	if _, err := saveProject(dir, p); err != nil {
		t.Fatalf("saveProject: %v", err)
	}
}

func TestUpdateDeployFileRemoveConfig(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("MINIPAAS_DEPLOY_VERSION", "0.0.1")
	defer os.Unsetenv("MINIPAAS_DEPLOY_VERSION")
	apps := filepath.Join(dir, "compose.apps.yaml")
	// project includes a config referenced by a service
	content := "version: '3.9'\nconfigs:\n  cfgA:\n    external: true\nservices:\n  api:\n    image: x\n    configs:\n      - source: cfgA\n        target: a\n"
	if err := os.WriteFile(apps, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := updateDeployFileRemoveConfig(dir, "cfgA"); err != nil {
		t.Fatalf("updateDeployFileRemoveConfig: %v", err)
	}
	// reload and assert removed
	p, _, err := loadProject(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if p.Configs != nil {
		t.Fatalf("top-level configs should be nil/empty: %#v", p.Configs)
	}
	if len(p.Services["api"].Configs) != 0 {
		t.Fatalf("service configs not cleared: %#v", p.Services["api"].Configs)
	}
}

func TestBuildDeployProject(t *testing.T) {
	p := buildDeployProject()
	if p == nil || p.Services == nil {
		t.Fatalf("buildDeployProject should initialize services map")
	}
}

func TestCreateDeployService_RequiresBuild(t *testing.T) {
	src := types.ServiceConfig{Image: "repo/app:1"}
	out := createDeployService(src, true)
	if out.Image != buildCommonImage(src.Image) {
		t.Fatalf("image not rewritten: %q", out.Image)
	}
	if _, ok := out.Networks["minipaas_network"]; !ok {
		t.Fatalf("network not set")
	}
	out2 := createDeployService(src, false)
	if out2.Image != "" {
		t.Fatalf("image should be empty when not requiring build: %q", out2.Image)
	}
}

func TestComposeEnsureDeployCreatesService(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	composeEnsureDeploy(p, "api")
	svc, ok := p.Services["api"]
	if !ok {
		t.Fatalf("service not created")
	}
	if svc.Networks == nil {
		t.Fatalf("networks not initialized")
	}
	if _, ok := svc.Networks["minipaas_network"]; !ok {
		t.Fatalf("minipaas_network missing: %#v", svc.Networks)
	}
}

func TestComposeEnsureDeployKeepsExisting(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["web"] = types.ServiceConfig{Networks: map[string]*types.ServiceNetworkConfig{"custom": {}}}
	composeEnsureDeploy(p, "web")
	if _, exists := p.Services["web"].Networks["custom"]; !exists {
		t.Fatalf("existing networks should be preserved")
	}
}

func TestAddComposeResilientDeploy(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["api"] = types.ServiceConfig{Image: "x"}
	if err := addComposeResilientDeploy(p, "api", "8080"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	svc := p.Services["api"]
	if svc.Deploy == nil || svc.Deploy.Mode != "replicated" || svc.Deploy.Replicas == nil || *svc.Deploy.Replicas != 2 {
		t.Fatalf("deploy not set as expected: %#v", svc.Deploy)
	}
	if svc.HealthCheck == nil || len(svc.HealthCheck.Test) == 0 {
		t.Fatalf("healthcheck not set: %#v", svc.HealthCheck)
	}
}

func TestAddComposeWorkerDeploy(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["worker"] = types.ServiceConfig{Image: "x"}
	if err := addComposeWorkerDeploy(p, "worker"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	svc := p.Services["worker"]
	if svc.Deploy == nil || svc.Deploy.Mode != "replicated" || svc.Deploy.Replicas == nil || *svc.Deploy.Replicas != 1 {
		t.Fatalf("deploy not set as expected: %#v", svc.Deploy)
	}
}

func TestAddComposeJobDeploy(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["job"] = types.ServiceConfig{Image: "x"}
	if err := addComposeJobDeploy(p, "job"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	svc := p.Services["job"]
	if svc.Deploy == nil || svc.Deploy.UpdateConfig == nil || svc.Deploy.UpdateConfig.Parallelism == nil || *svc.Deploy.UpdateConfig.Parallelism != 0 {
		t.Fatalf("job update config not set: %#v", svc.Deploy)
	}
	if svc.Deploy.RestartPolicy == nil || svc.Deploy.RestartPolicy.Condition != "on-failure" {
		t.Fatalf("job restart policy not set: %#v", svc.Deploy)
	}
}

func TestAddComposeCronDeploy(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["cron"] = types.ServiceConfig{Image: "x"}
	if err := addComposeCronDeploy(p, "cron", "* * * * *"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	svc := p.Services["cron"]
	if svc.Deploy == nil || svc.Deploy.Replicas == nil || *svc.Deploy.Replicas != 0 {
		t.Fatalf("cron replicas not set to 0: %#v", svc.Deploy)
	}
	if svc.Deploy.Labels["swarm.cronjob.enable"] != "true" {
		t.Fatalf("cron label missing: %#v", svc.Deploy.Labels)
	}
}

func TestAddComposeSecretAndConfig_NoDupes(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["api"] = types.ServiceConfig{Image: "x"}

	if err := addComposeSecret(p, "secret-hash", "env", []string{"api"}); err != nil {
		t.Fatalf("add secret: %v", err)
	}
	if err := addComposeSecret(p, "secret-hash", "env", []string{"api"}); err != nil {
		t.Fatalf("add secret twice: %v", err)
	}
	if len(p.Secrets) != 1 || len(p.Services["api"].Secrets) != 1 {
		t.Fatalf("secret dupes detected: secrets=%d svcRefs=%d", len(p.Secrets), len(p.Services["api"].Secrets))
	}

	if err := addComposeConfig(p, "cfg-hash", "app.conf", []string{"api"}); err != nil {
		t.Fatalf("add config: %v", err)
	}
	if err := addComposeConfig(p, "cfg-hash", "app.conf", []string{"api"}); err != nil {
		t.Fatalf("add config twice: %v", err)
	}
	if len(p.Configs) != 1 || len(p.Services["api"].Configs) != 1 {
		t.Fatalf("config dupes detected: configs=%d svcRefs=%d", len(p.Configs), len(p.Services["api"].Configs))
	}
}

func TestAddComposeSecretMissingService(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	if err := addComposeSecret(p, "sec", "file", []string{"missing"}); err == nil {
		t.Fatalf("expected error for missing service")
	}
}

func TestAddComposeConfigMissingService(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	if err := addComposeConfig(p, "cfg", "file", []string{"missing"}); err == nil {
		t.Fatalf("expected error for missing service")
	}
}

func TestRemoveConfigFromDeployProject(t *testing.T) {
	p := &types.Project{Services: make(types.Services), Configs: map[string]types.ConfigObjConfig{
		"cfgA": {External: true},
		"cfgB": {External: true},
	}}
	p.Services["api"] = types.ServiceConfig{
		Name: "api",
		Configs: []types.ServiceConfigObjConfig{
			{Source: "cfgA", Target: "a"},
			{Source: "cfgB", Target: "b"},
		},
	}

	changed := removeConfigFromDeployProject(p, "cfgA")
	if !changed {
		t.Fatalf("expected changes true")
	}
	if _, ok := p.Configs["cfgA"]; ok {
		t.Fatalf("cfgA should be removed from top-level")
	}
	if len(p.Services["api"].Configs) != 1 || p.Services["api"].Configs[0].Source != "cfgB" {
		t.Fatalf("service configs not filtered: %#v", p.Services["api"].Configs)
	}
}

func TestBuildCommonImage(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"myapp:latest", "registry:5000/myapp:${MINIPAAS_DEPLOY_VERSION}"},
		{"myapp", "registry:5000/myapp:${MINIPAAS_DEPLOY_VERSION}"},
		{"repo/myapp:1.2.3", "repo/myapp:${MINIPAAS_DEPLOY_VERSION}"},
		{"registry:5001/repo/myapp:abc", "registry:5001/repo/myapp:${MINIPAAS_DEPLOY_VERSION}"},
	}
	for _, tt := range tests {
		got := buildCommonImage(tt.in)
		if got != tt.want {
			t.Fatalf("buildCommonImage(%q)=%q want %q", tt.in, got, tt.want)
		}
	}
}
