package main

import (
	"context"
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
)

func composeLoadProject(files []string) (*types.Project, error) {
	ctx := context.Background()

	options, err := cli.NewProjectOptions(
		files,
		cli.WithName("minipaas"),
		cli.WithResolvedPaths(false),
		cli.WithConsistency(false),
		cli.WithoutEnvironmentResolution,
	)
	if err != nil {
		return nil, err
	}

	project, err := options.LoadProject(ctx)
	if err != nil {
		return nil, err
	}

	return project, err
}

func composeLoadDeployProject(files []string) (*types.Project, error) {
	ctx := context.Background()

	options, err := cli.NewProjectOptions(
		files,
		cli.WithName("minipaas"),
		cli.WithResolvedPaths(false),
		cli.WithConsistency(false),
		cli.WithOsEnv,
	)
	if err != nil {
		return nil, err
	}

	project, err := options.LoadProject(ctx)
	if err != nil {
		return nil, err
	}

	return project, err
}
