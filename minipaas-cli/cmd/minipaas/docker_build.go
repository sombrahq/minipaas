package main

import (
	"fmt"
	"path/filepath"

	"github.com/compose-spec/compose-go/v2/types"
)

func buildCommandFromService(svc types.ServiceConfig) []string {
	cmd := []string{"docker", "build", "--network", "host"}

	// If an image name is specified, use it to tag the built image.
	if svc.Image != "" {
		cmd = append(cmd, "-t", svc.Image)
	}

	// If the service has build configuration, process it.
	if svc.Build != nil {
		context := svc.Build.Context
		if context == "" {
			context = "."
		}

		// Specify the Dockerfile if provided.
		if svc.Build.Dockerfile != "" {
			cmd = append(cmd, "-f", filepath.Join(context, svc.Build.Dockerfile))
		}

		// Append build arguments if any.f
		if svc.Build.Args != nil {
			for key, valPtr := range svc.Build.Args {
				if valPtr != nil {
					cmd = append(cmd, "--build-arg", fmt.Sprintf("%s=%s", key, *valPtr))
				} else {
					cmd = append(cmd, "--build-arg", fmt.Sprintf("%s=", key))
				}
			}
		}

		// Use the provided build context. If empty, default to current directory.
		cmd = append(cmd, svc.Build.Context)
	} else {
		// If no build configuration is provided, default to current directory.
		cmd = append(cmd, ".")
	}

	return cmd
}
