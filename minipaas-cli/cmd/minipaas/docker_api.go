package main

// DockerService represents the minimal structure from "docker service inspect"
// needed to extract the secrets used by a service.
type DockerService struct {
	Spec struct {
		TaskTemplate struct {
			ContainerSpec struct {
				Secrets []struct {
					SecretName string `json:"SecretName"`
				} `json:"Secrets"`
				Configs []struct {
					ConfigName string `json:"ConfigName"`
				} `json:"Configs"`
			} `json:"ContainerSpec"`
		} `json:"TaskTemplate"`
	} `json:"Spec"`
}
