package cmd

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"io.twasyl/devcore/pkg/config"
)

type devcoreContainer struct {
	testcontainers.Container
}

func createDevcoreContainer(ctx context.Context) (*devcoreContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:      "alpine",
		AutoRemove: true,
		Cmd:        []string{"sleep", "infinity"},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	if err := container.CopyFileToContainer(ctx, "../devcore", "/usr/local/bin/devcore", 0755); err != nil {
		return nil, err
	}

	if code, err := container.Exec(ctx, []string{"devcore", "version"}); err != nil {
		return nil, err
	} else if code != 0 {
		return nil, errors.New(fmt.Sprintf("`devcore version` returned %d", code))
	}

	return &devcoreContainer{container}, nil
}

func TestToolsInstallation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	container, err := createDevcoreContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	for _, tool := range config.DefaultSupportedTools {
		if exitCode, err := container.Exec(ctx, []string{"devcore", "tools", "install", tool.Name}); err != nil {
			t.Errorf("Error installing %s: %s", tool.Name, err)
		} else if exitCode != 0 {
			t.Errorf("Installing %s returned %d", tool.Name, exitCode)
		} else if whichExitCode, err := container.Exec(ctx, []string{"which", tool.CommandLineName}); err != nil {
			t.Errorf("Error executing which for tool %s", tool.Name)
		} else if whichExitCode != 0 {
			t.Errorf("%s is not installed", tool.Name)
		}
	}
}
