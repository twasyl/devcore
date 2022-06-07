package config

import "fmt"

type DockerComposeContextNotFound struct {
	Name string
}

func (e *DockerComposeContextNotFound) Error() string {
	return fmt.Sprintf("Docker compose context not found: '%s'", e.Name)
}

func IsDockerComposeContextNotFound(err error) bool {
	if err != nil {
		_, yes := err.(*DockerComposeContextNotFound)
		return yes
	}
	return false
}

type JenkinsContextNotFound struct {
	Name string
}

func (e *JenkinsContextNotFound) Error() string {
	return fmt.Sprintf("Jenkins context not found: '%s'", e.Name)
}

func IsJenkinsContextNotFound(err error) bool {
	if err != nil {
		_, yes := err.(*JenkinsContextNotFound)
		return yes
	}
	return false
}

type ToolNotFound struct {
	Name string
}

func (e *ToolNotFound) Error() string {
	return fmt.Sprintf("Tool not found: '%s'", e.Name)
}

func IsToolNotFound(err error) bool {
	if err != nil {
		_, yes := err.(*ToolNotFound)
		return yes
	}
	return false
}

type ServerNotFound struct {
	Name string
}

func (e *ServerNotFound) Error() string {
	return fmt.Sprintf("Server not found: '%s'", e.Name)
}

func IsServerNotFound(err error) bool {
	if err != nil {
		_, yes := err.(*ServerNotFound)
		return yes
	}
	return false
}
